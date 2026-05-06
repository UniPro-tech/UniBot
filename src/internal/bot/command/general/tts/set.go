package tts

import (
	"unibot/internal/bot/command/general/tts/ttsSet"

	"github.com/disgoorg/disgo/discord"
)

func LoadSetCommandContext() discord.ApplicationCommandOptionSubCommandGroup {
	return discord.ApplicationCommandOptionSubCommandGroup{
		Name:        "set",
		Description: "TTSの設定を変更します",
		Options: []discord.ApplicationCommandOptionSubCommand{
			ttsSet.LoadVoiceCommandContext(),
			ttsSet.LoadSpeedCommandContext(),
		},
	}
}
