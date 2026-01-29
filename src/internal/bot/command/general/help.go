package general

import (
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/bot/command/general/help"

	"github.com/bwmarrin/discordgo"
)

func LoadHelpCommandContext() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "利用可能なコマンドの一覧を表示します",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "command",
				Description: "特定のコマンドのヘルプを表示します",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "About",
						Value: "/about",
					},
					{
						Name:  "Help",
						Value: "/help",
					},
					{
						Name:  "Ping",
						Value: "/ping",
					},
					{
						Name:  "TTS",
						Value: "/tts <subcommand>",
					},
					{
						Name:  "TTS Dict",
						Value: "/tts dict <subcommands>",
					},
					{
						Name:  "Skip",
						Value: "/skip",
					},
				},
			},
		},
	}
}

func Help(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
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
						Description: "ヘルプの表示に失敗しました。",
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
	user := i.User
	if i.Member != nil {
		user = i.Member.User
	}

	// コマンドオプションを取得
	options := i.ApplicationCommandData().Options
	var specificCommand string
	if len(options) > 0 {
		specificCommand = options[0].StringValue()
	}

	var fields []*discordgo.MessageEmbedField
	var title string
	var description string

	// 特定のコマンドが指定された場合
	if specificCommand != "" {
		title = "Help - " + specificCommand + " コマンド"
		description = "コマンドの詳細情報です。"

		// 指定されたコマンドのフィールドのみを抽出
		for _, cmd := range help.HelpCommands {
			if cmd.Name == specificCommand {
				fields = append(fields, &discordgo.MessageEmbedField{
					Name:   "説明",
					Value:  cmd.Description,
					Inline: false,
				})

				// Usage が設定されている場合は表示
				if cmd.Usage != "" {
					fields = append(fields, &discordgo.MessageEmbedField{
						Name:   "使い方",
						Value:  cmd.Usage,
						Inline: false,
					})
				}
				break
			}
		}

		// コマンドが見つからない場合
		if len(fields) == 0 {
			description = "指定されたコマンドが見つかりませんでした。"
		}
	} else {
		// すべてのコマンドを表示
		title = "Help - 利用可能なコマンド一覧"
		description = "以下は利用可能なコマンドの一覧です。"

		fields = make([]*discordgo.MessageEmbedField, 0, len(help.HelpCommands))
		for _, cmd := range help.HelpCommands {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  cmd.Name,
				Value: cmd.Description,
			})
		}
	}

	responseEmbed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       config.Colors.Primary,
		Fields:      fields,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: user.AvatarURL(""),
			Name:    user.Username,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Requested by " + user.Username,
			IconURL: user.AvatarURL(""),
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
					Description: "ヘルプを表示できませんでした。",
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
