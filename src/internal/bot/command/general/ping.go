package general

import (
	"time"
	"unibot/internal"

	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := internal.LoadConfig()

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
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{responseEmbed},
		},
	})
}
