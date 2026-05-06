package ttsSet

import (
	"fmt"
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
	}
}

func Voice(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config

		speakers, err := ttsutil.FetchSpeakers(ctx)
		if err != nil {
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

		currentSpeakerID := ttsutil.GetCurrentSpeakerID(ctx, e.User().ID.String())
		content, components := ttsutil.BuildVoiceMessage(0, pages, currentSpeakerID)

		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithContent(content).WithComponents(components...))
		return err
	}
}
