package general

import (
	"time"
	"unibot/internal"
	"unibot/internal/bot/command/general/tts"
	"unibot/internal/bot/command/general/tts/dict"

	"github.com/bwmarrin/discordgo"
)

func LoadTtsCommandContext() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "tts",
		Description: "テキスト読み上げを行います",
		Options: []*discordgo.ApplicationCommandOption{
			tts.LoadJoinCommandContext(),
			tts.LoadLeaveCommandContext(),
			tts.LoadSkipCommandContext(),
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "dict",
				Description: "TTS辞書を管理します",
				Options: []*discordgo.ApplicationCommandOption{
					dict.LoadAddCommandContext(),
					dict.LoadRemoveCommandContext(),
					dict.LoadListCommandContext(),
				},
			},
		},
	}
}

var ttsHandler = map[string]func(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate){
	"join":  tts.Join,
	"leave": tts.Leave,
	"skip":  tts.Skip,
}

var ttsDictHandler = map[string]func(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate){
	"add":    dict.Add,
	"remove": dict.Remove,
	"list":   dict.List,
}

func Tts(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	subCommandGroup := i.ApplicationCommandData().Options[0]

	// サブコマンドグループの場合
	if subCommandGroup.Type == discordgo.ApplicationCommandOptionSubCommandGroup {
		if subCommandGroup.Name == "dict" {
			subCommand := subCommandGroup.Options[0]
			if handler, exists := ttsDictHandler[subCommand.Name]; exists {
				handler(ctx, s, i)
				return
			}
		}
	} else {
		// サブコマンドの場合
		if handler, exists := ttsHandler[subCommandGroup.Name]; exists {
			handler(ctx, s, i)
			return
		}
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "不明なサブコマンドです。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}
