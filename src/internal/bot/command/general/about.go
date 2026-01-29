package general

import (
	"log"
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
						Description: "情報の表示に失敗しました。",
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

	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{responseEmbed},
	})
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "情報を表示できませんでした。",
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
