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
		},
	}
}

var ttsHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"join": tts.Join,
}

func Tts(s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := internal.LoadConfig()
	subCommandGroup := i.ApplicationCommandData().Options[0]
	if subCommandGroup.Type == discordgo.ApplicationCommandOptionSubCommandGroup {
		if handler, exists := ttsHandler[subCommandGroup.Name]; exists {
			handler(s, i)
			return
		}
	} else {
		if handler, exists := ttsHandler[subCommandGroup.Name]; exists {
			handler(s, i)
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
