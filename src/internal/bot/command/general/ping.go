package general

import (
	"log"
	"time"
	"unibot/internal"

	"github.com/bwmarrin/discordgo"
)

func LoadPingCommandContext() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "スピードテストを行います",
	}
}

func Ping(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config

	// Get ws websocketLatency
	websocketLatency := s.HeartbeatLatency()

	// Respond to the interaction
	responseEmbed := &discordgo.MessageEmbed{
		Title:       "Pong 🏓",
		Description: "スピードテストの結果です",
		Color:       config.Colors.Primary,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "WebSocket Latency",
				Value: websocketLatency.String(),
			},
		},
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: i.Member.AvatarURL(""),
			Name:    i.Member.DisplayName(),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Requested by " + i.Member.DisplayName(),
			IconURL: i.Member.AvatarURL(""),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	done := make(chan struct{})
	go func() {
		select {
		case <-done:
			return
		case <-time.After(3 * time.Minute):
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "スピードテストに失敗しました。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			})
			if err != nil {
				log.Println("Failed to edit deferred interaction on timeout:", err)
			}
		}
	}()
	defer close(done)

	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{responseEmbed},
	})
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "スピードテストに失敗しました。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}
}
