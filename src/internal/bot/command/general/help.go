package general

import (
	"time"
	"unibot/internal"
	"unibot/internal/bot/command/general/help"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func LoadHelpCommandContext() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "help",
		Description: "利用可能なコマンドの一覧を表示します",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "command",
				Description: "特定のコマンドのヘルプを表示します",
				Required:    false,
				Choices: []discord.ApplicationCommandOptionChoiceString{
					{
						Name:  "About",
						Value: "/about",
					},
					{
						Name:  "ColorCode",
						Value: "/colorcode",
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
						Name:  "TTS Skip",
						Value: "/tts skip",
					},
					{
						Name:  "TTS Speed",
						Value: "/tts set speed <value>",
					},
				},
			},
		},
	}
}

func Help(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config

		// コマンドオプションを取得
		specificCommand, exist := data.OptString("command")

		var fields []discord.EmbedField
		var title string
		var description string

		// 特定のコマンドが指定された場合
		if exist {
			title = "Help - " + specificCommand + " コマンド"
			description = "コマンドの詳細情報です。"

			// 指定されたコマンドのフィールドのみを抽出
			for _, cmd := range help.HelpCommands {
				if cmd.Name == specificCommand {
					fields = append(fields, discord.EmbedField{
						Name:  "説明",
						Value: cmd.Description,
					})

					// Usage が設定されている場合は表示
					if cmd.Usage != "" {
						fields = append(fields, discord.EmbedField{
							Name:  "使い方",
							Value: cmd.Usage,
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

			fields = make([]discord.EmbedField, 0, len(help.HelpCommands))
			for _, cmd := range help.HelpCommands {
				fields = append(fields, discord.EmbedField{
					Name:  cmd.Name,
					Value: cmd.Description,
				})
			}
		}

		responseEmbed := discord.Embed{
			Title:       title,
			Description: description,
			Color:       config.Colors.Primary,
			Fields:      fields,
			Footer: &discord.EmbedFooter{
				Text:    "Requested by " + e.User().Username,
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}

		_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))
		return err
	}
}
