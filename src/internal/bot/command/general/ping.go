package general

import (
	"fmt"
	"time"
	"unibot/internal"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func LoadPingCommandContext() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "ping",
		Description: "スピードテストを行います",
	}
}

func Ping(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config

		// Get ws websocketLatency
		websocketLatency := e.Client().Gateway.Latency()

		// Respond to the interaction
		responseEmbed := discord.Embed{
			Title:       "Pong 🏓",
			Description: "スピードテストの結果です",
			Color:       config.Colors.Primary,
			Fields: []discord.EmbedField{
				{
					Name:  "WebSocket Latency",
					Value: websocketLatency.String(),
				},
			},
			Footer: &discord.EmbedFooter{
				Text:    fmt.Sprintf("Requested by %s", e.User().Username),
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}

		_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(),
			discord.NewMessageCreate().WithEmbeds(responseEmbed))
		return err
	}
}
