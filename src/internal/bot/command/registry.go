package command

import (
	"unibot/internal"
	"unibot/internal/bot/command/admin"
	"unibot/internal/bot/command/general"
	"unibot/internal/bot/command/server_management"

	"github.com/bwmarrin/discordgo"
)

type HandlerEntry struct {
	Handler   func(*internal.BotContext, *discordgo.Session, *discordgo.InteractionCreate)
	Ephemeral bool
}

var Handlers = map[string]HandlerEntry{
	"ping": {
		Handler: general.Ping,
	},
	"about": {
		Handler: general.About,
	},
	"maintenance": {
		Handler:   admin.Maintenance,
		Ephemeral: true,
	},
	"tts": {
		Handler: general.Tts,
	},
	"help": {
		Handler: general.Help,
	},
	"colorcode": {
		Handler: general.ColorCode,
	},
	"rolepanel": {
		Handler: server_management.Rolepanel,
	},
}
