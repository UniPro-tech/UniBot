package command

import (
	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get ws websocketLatency
	websocketLatency := s.HeartbeatLatency()

	// Respond to the interaction
	responseEmbed := &discordgo.MessageEmbed{
		Title:       "Pong 🏓",
		Description: "スピードテストの結果です",
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
			Text:    "UniBot",
			IconURL: s.State.User.AvatarURL(""),
		},
	}
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{responseEmbed},
		},
	})
}
