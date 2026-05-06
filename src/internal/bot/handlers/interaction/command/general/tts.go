package general

import (
	"unibot/internal/bot/handlers/interaction/command/general/tts"

	"github.com/disgoorg/disgo/discord"
)

func LoadTtsCommandContext() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "tts",
		Description: "テキスト読み上げを行います",
		Options: []discord.ApplicationCommandOption{
			tts.LoadJoinCommandContext(),
			tts.LoadLeaveCommandContext(),
			tts.LoadSkipCommandContext(),
			tts.LoadDictCommandContext(),
			tts.LoadSetCommandContext(),
		},
	}
}
