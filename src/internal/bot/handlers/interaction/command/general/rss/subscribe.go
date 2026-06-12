package rss

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	neturl "net/url"
	"sort"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/mmcdole/gofeed"
)

func LoadSubscribeCommandContext() discord.ApplicationCommandOption {
	return discord.ApplicationCommandOptionSubCommand{
		Name:        "subscribe",
		Description: "RSSフィードを購読します",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "url",
				Description: "RSSフィードのURLを設定します",
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				Name:        "title",
				Description: "タイトルを設定します",
				Required:    false,
			},
		},
	}
}

func Subscribe(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config

		if e.Channel().Type() == discord.ChannelTypeDM || e.Channel().Type() == discord.ChannelTypeGroupDM {
			responseEmbed := discord.Embed{
				Title:       "DMでは実行できません",
				Description: "このコマンドはDMでは実行できません。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed).WithEphemeral(true))
			return err
		}

		guildID := *e.GuildID()
		var title *string
		var url string
		for _, opt := range data.Options {
			switch opt.Name {
			case "title":
				title = func() *string {
					titleValue := opt.String()
					return &titleValue
				}()
			case "url":
				url = opt.String()
			}
		}

		const (
			maxFeedSize = 10 * 1024 * 1024 // 10MB
			maxRedirect = 5
		)

		parsedURL, err := neturl.Parse(url)
		if err != nil {
			return errorSubscribeResponse(config, e)
		}

		if err := validateURL(parsedURL); err != nil {
			return errorSubscribeResponse(config, e)
		}

		dialer := &net.Dialer{}

		transport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}

				ips, err := net.LookupIP(host)
				if err != nil {
					return nil, err
				}

				if len(ips) == 0 {
					return nil, errors.New("no ip found")
				}

				// 全IP検査
				for _, ip := range ips {
					if isPrivateIP(ip) {
						return nil, fmt.Errorf("private ip detected: %s", ip)
					}
				}

				// 検査済みIPへ接続
				return dialer.DialContext(
					ctx,
					network,
					net.JoinHostPort(ips[0].String(), port),
				)
			},
		}

		httpClient := &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= maxRedirect {
					return errors.New("too many redirects")
				}

				return validateURL(req.URL)
			},
		}
		resp, err := httpClient.Get(url)
		if err != nil {
			return errorSubscribeResponse(config, e)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return errorSubscribeResponse(config, e)
		}
		limitedBody := io.LimitReader(resp.Body, 10<<20)
		fp := gofeed.NewParser()
		feed, err := fp.Parse(limitedBody)
		if err != nil {
			return errorSubscribeResponse(config, e)
		}

		if feed.Title != "" && title == nil {
			title = &feed.Title
		}
		var hash *string
		if len(feed.Items) != 0 {
			// 新しい日時がindex:0
			sort.SliceStable(feed.Items, func(i, j int) bool {
				a := feed.Items[i].PublishedParsed
				b := feed.Items[j].PublishedParsed

				switch {
				case a == nil && b == nil:
					return false
				case a == nil:
					return false
				case b == nil:
					return true
				default:
					return a.After(*b)
				}
			})
			hash = func() *string {
				data := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", feed.Items[0].Title, feed.Items[0].Description)))
				return &data
			}()
		}

		db := ctx.DB
		guildRepo := repository.NewGuildRepository(db)
		if _, err := guildRepo.GetOrCreate(guildID.String()); err != nil {
			return err
		}
		rssRepo := repository.NewRSSSettingRepository(db)
		if err := rssRepo.Create(&model.RSSSetting{
			GuildID:                      guildID.String(),
			ChannelID:                    e.Channel().ID().String(),
			URL:                          url,
			Title:                        title,
			LastItemTitleDescriptionHash: hash,
		}); err != nil {
			return err
		}

		// 成功レスポンス
		responseEmbed := discord.Embed{
			Title:       "RSS購読",
			Description: "RSS購読設定が完了しました。",
			Fields: []discord.EmbedField{
				{
					Name:  "URL",
					Value: url,
				},
			},
			Color: config.Colors.Success,
			Footer: &discord.EmbedFooter{
				Text:    fmt.Sprintf("Requested by %s", e.User().Username),
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}
		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed).WithEphemeral(false))
		return err
	}
}

func isPrivateIP(ip net.IP) bool {
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	}

	privateCIDRs := []string{
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",

		"100.64.0.0/10",
		"198.18.0.0/15",
		"224.0.0.0/4",
		"240.0.0.0/4",

		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}

	for _, cidr := range privateCIDRs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}

		if network.Contains(ip) {
			return true
		}
	}

	return false
}

func errorSubscribeResponse(config *internal.Config, e *handler.CommandEvent) error {
	responseEmbed := discord.Embed{
		Title:       "RSS購読",
		Description: "RSSフィードの取得に失敗しました。",
		Color:       config.Colors.Error,
		Footer: &discord.EmbedFooter{
			Text:    fmt.Sprintf("Requested by %s", e.User().Username),
			IconURL: e.User().EffectiveAvatarURL(),
		},
		Timestamp: func() *time.Time {
			t := time.Now()
			return &t
		}(),
	}
	_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed).WithEphemeral(true))
	return err
}

func validateURL(u *url.URL) error {
	if u == nil {
		return errors.New("nil url")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("invalid scheme")
	}

	host := u.Hostname()
	if host == "" {
		return errors.New("missing host")
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return err
	}

	if len(ips) == 0 {
		return errors.New("no ip found")
	}

	for _, ip := range ips {
		if isPrivateIP(ip) {
			return fmt.Errorf("private address: %s", ip)
		}
	}

	port := u.Port()
	if port != "" && port != "80" && port != "443" {
		return fmt.Errorf("invalid port: %s", port)
	}

	return nil
}
