package tts

import (
	"unibot/internal"
	"unibot/internal/bot/command/general/tts/set"

	"github.com/bwmarrin/discordgo"
)

func LoadSetCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
		Name:        "set",
		Description: "TTSの設定を変更します",
		Options: []*discordgo.ApplicationCommandOption{
			set.LoadVoiceCommandContext(),
		},
	}
}

var setHandler = map[string]func(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate){
	"voice": set.Voice,
}

func Set(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	subCommandGroup := i.ApplicationCommandData().Options[0]
	subCommand := subCommandGroup.Options[0]

	if handler, exists := setHandler[subCommand.Name]; exists {
		handler(ctx, s, i)
	}
}
