package general

import (
	"time"
	"unibot/internal"
	"unibot/internal/bot/command/general/tts"

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
			tts.LoadDictCommandContext(),
			tts.LoadSetCommandContext(),
			tts.LoadSpeedCommandContext(),
		},
	}
}

var ttsHandler = map[string]func(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate){
	"join":  tts.Join,
	"leave": tts.Leave,
	"skip":  tts.Skip,
	"dict":  tts.Dict,
	"set":   tts.Set,
	"speed": tts.Speed,
}

func Tts(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	subCommandGroup := i.ApplicationCommandData().Options[0]

	// サブコマンドグループの場合
	if subCommandGroup.Type == discordgo.ApplicationCommandOptionSubCommandGroup {
		if handler, exists := ttsHandler[subCommandGroup.Name]; exists {
			handler(ctx, s, i)
			return
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
