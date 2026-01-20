package command

import (
	"github.com/bwmarrin/discordgo"

	"unibot/internal"
)

func About(s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := internal.LoadConfig()
	// コントリビューター一覧をMarkdown形式で作成
	contributorsText := ""
	// Botは最後にする
	for _, contributor := range config.Contributors {
		if contributor.IsBot {
			continue
		}
		contributorsText += "- [" + contributor.Username + "](" + contributor.Profile + ")\n"
	}
	for _, contributor := range config.Contributors {
		if !contributor.IsBot {
			continue
		}
		contributorsText += "- [" + contributor.Username + "](" + contributor.Profile + ")\n"
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
				Name:  "コントリビューター",
				Value: contributorsText,
			},
			{
				Name:  "GitHub",
				Value: config.GitHub,
			},
			{
				Name:  "サポートサーバー",
				Value: config.SupportServer,
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
