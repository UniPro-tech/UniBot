package command

import (
	"github.com/bwmarrin/discordgo"

	"unibot/internal"
)

func About(s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := internal.LoadConfig()
	contributorsText := ""
	for _, contributor := range config.Contributors {
		contributorsText += "- " + contributor + "\n"
	}

	responseEmbed := &discordgo.MessageEmbed{
		Title:       "About UniBot 🤖",
		Description: "UniBotはデジタル創作サークルUniProjectの内製Botです。",
		Color:       config.Colors.Primary,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "バージョン",
				Value: config.BotVersion,
			},
			{
				Name:  "開発者",
				Value: contributorsText,
			},
			{
				Name:  "GitHub",
				Value: config.GitHub,
			},
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{responseEmbed},
		},
	})
}
