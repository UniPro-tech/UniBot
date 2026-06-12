package rss

import (
	"encoding/base64"
	"fmt"
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

		// 初回Fetch
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(url)
		if err != nil {
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

		if feed.Title != "" && title == nil {
			title = &feed.Title
		}
		// 新しい日時がindex:0
		sort.Slice(feed.Items, func(i, j int) bool {
			prev := feed.Items[i]
			next := feed.Items[j]
			if prev.PublishedParsed != nil && next.PublishedParsed != nil {
				prevNano := prev.PublishedParsed.UnixNano()
				nextNano := next.PublishedParsed.UnixNano()
				return prevNano >= nextNano
			}
			return false
		})
		hash := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", feed.Items[0].Title, feed.Items[0].Description)))

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
			LastItemTitleDescriptionHash: &hash,
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
