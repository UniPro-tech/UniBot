package ttsSet

import (
	"fmt"
	"log/slog"
	"time"
	"unibot/internal"
	"unibot/internal/bot/ttsutil"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func LoadVoiceCommandContext() discord.ApplicationCommandOptionSubCommand {
	return discord.ApplicationCommandOptionSubCommand{
		Name:        "voice",
		Description: "読み上げの話者を設定します",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionBool{
				Name:        "global",
				Description: "グローバルに設定するか (デフォルト: false)",
				Required:    false,
			},
		},
	}
}

func Voice(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config

		if _, ok := e.Guild(); !ok {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "DMで実行することはできません。",
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

		speakers, err := ttsutil.FetchSpeakers(ctx)
		if err != nil {
			slog.ErrorContext(e.Ctx, "failed to fetch tts speaker list", slog.Any("err", err))
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "話者情報の取得に失敗しました。",
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

		pages := ttsutil.BuildSpeakerPages(speakers, ttsutil.SpeakerPageSize)
		if len(pages) == 0 {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "話者情報が見つかりませんでした。",
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

		global := false
		if value, ok := data.Options["global"]; ok {
			global = value.Bool()
		}

		currentSpeakerID, err := ttsutil.GetCurrentSpeakerID(ctx, e.User().ID, global, e.GuildID())
		if err != nil {
			slog.ErrorContext(e.Ctx, "failed to fetch tts settings", slog.Any("err", err))
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "設定の取得に失敗しました。",
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
		embed, components, err := ttsutil.BuildVoiceMessage(ctx, 0, pages, currentSpeakerID, global)
		if err != nil {
			slog.ErrorContext(e.Ctx, "failed to generate response message", slog.Any("err", err))
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "メッセージの生成に失敗しました。",
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

		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*embed).WithComponents(components...))
		return err
	}
}
