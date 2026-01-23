package tts

import (
	"unibot/internal"
	"unibot/internal/bot/command/general/tts/dict"

	"github.com/bwmarrin/discordgo"
)

func LoadDictCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
		Name:        "dict",
		Description: "TTS辞書を管理します",
		Options: []*discordgo.ApplicationCommandOption{
			dict.LoadAddCommandContext(),
			dict.LoadRemoveCommandContext(),
			dict.LoadListCommandContext(),
		},
	}
}

var dictHandler = map[string]func(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate){
	"add":    dict.Add,
	"remove": dict.Remove,
	"list":   dict.List,
}

func Dict(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	subCommandGroup := i.ApplicationCommandData().Options[0]
	subCommand := subCommandGroup.Options[0]

	if handler, exists := dictHandler[subCommand.Name]; exists {
		handler(ctx, s, i)
	}
}
