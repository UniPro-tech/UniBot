package general

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"unibot/internal"
)

func LoadAboutCommandContext() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "about",
		Description: "ボットの情報を表示します",
	}
}

func About(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
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
		Title:       "About " + config.BotName + " 🤖",
		Description: config.Description,
		Color:       config.Colors.Primary,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Version",
				Value: config.BotVersion,
			},
			{
				Name:  "Contributors",
				Value: contributorsText,
			},
			{
				Name:  "GitHub",
				Value: config.GitHub,
			},
			{
				Name:  "Support Server",
				Value: config.SupportServer,
			},
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
