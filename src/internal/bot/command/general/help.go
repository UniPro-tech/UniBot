package general

import (
	"encoding/json"
	"time"
	"unibot/internal"

	_ "embed"

	"github.com/bwmarrin/discordgo"
)

type helpCommandData struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Fields      []helpCommandField `json:"fields"`
}

type helpCommandField struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	HowToUse string `json:"How to use"`
}

//go:embed help/commands.json
var helpCommandsJSON []byte

func loadHelpCommands() helpCommandData {
	var data helpCommandData
	if err := json.Unmarshal(helpCommandsJSON, &data); err != nil {
		return helpCommandData{
			Title:       "Help - 利用可能なコマンド一覧",
			Description: "コマンド情報の読み込みに失敗しました。",
		}
	}

	return data
}

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

	commandData := loadHelpCommands()

	var fields []*discordgo.MessageEmbedField
	var title string
	var description string

	// 特定のコマンドが指定された場合
	if specificCommand != "" {
		title = "Help - " + specificCommand + " コマンド"
		description = "コマンドの詳細情報です。"

		// 指定されたコマンドのフィールドのみを抽出
		for _, field := range commandData.Fields {
			if field.Name == specificCommand {
				fields = append(fields, &discordgo.MessageEmbedField{
					Name:   "説明",
					Value:  field.Value,
					Inline: false,
				})

				// How to use が設定されている場合は表示
				if field.HowToUse != "" {
					fields = append(fields, &discordgo.MessageEmbedField{
						Name:   "使い方",
						Value:  field.HowToUse,
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

		fields = make([]*discordgo.MessageEmbedField, 0, len(commandData.Fields))
		for _, field := range commandData.Fields {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  field.Name,
				Value: field.Value,
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

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{responseEmbed},
		},
	})
}
