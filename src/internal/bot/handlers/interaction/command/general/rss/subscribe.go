package rss

import (
	"fmt"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
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

		db := ctx.DB
		guildRepo := repository.NewGuildRepository(db)
		guildRepo.GetOrCreate(guildID.String())
		rssRepo := repository.NewRSSSettingRepository(db)
		rssRepo.Create(&model.RSSSetting{
			GuildID: guildID.String(),
			URL:     url,
			Title:   title,
		})

		// 成功レスポンス
		responseEmbed := discord.Embed{
			Title:       "TTSボイスチャンネル接続",
			Description: "ボイスチャンネルに参加しました。",
			Color:       config.Colors.Success,
			Footer: &discord.EmbedFooter{
				Text:    fmt.Sprintf("Requested by %s", e.User().Username),
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}
		_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed).WithEphemeral(false))
		return err
	}
}
