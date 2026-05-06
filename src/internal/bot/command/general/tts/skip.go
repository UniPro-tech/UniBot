package tts

import (
	"fmt"
	"time"
	"unibot/internal"
	"unibot/internal/bot/voice"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func LoadSkipCommandContext() discord.ApplicationCommandOptionSubCommand {
	return discord.ApplicationCommandOptionSubCommand{
		Name:        "skip",
		Description: "現在再生中の音声をスキップします",
	}
}

func Skip(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config

		player := voice.GetManager().Get(e.GuildID().String())
		if player == nil {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "現在再生中の音声はありません。",
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

		player.SkipCurrent()

		responseEmbed := discord.Embed{
			Title:       "音声スキップ",
			Description: "現在再生中の音声をスキップしました。",
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
