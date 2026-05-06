package tts

import (
	"unibot/internal/bot/handlers/interaction/command/general/tts/dict"

	"github.com/disgoorg/disgo/discord"
)

func LoadDictCommandContext() discord.ApplicationCommandOptionSubCommandGroup {
	return discord.ApplicationCommandOptionSubCommandGroup{
		Name:        "dict",
		Description: "TTS辞書を管理します",
		Options: []discord.ApplicationCommandOptionSubCommand{
			dict.LoadAddCommandContext(),
			dict.LoadRemoveCommandContext(),
			dict.LoadListCommandContext(),
		},
	}
}
