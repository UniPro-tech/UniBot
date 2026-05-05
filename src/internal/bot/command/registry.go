package command

import (
	"unibot/internal"
	"unibot/internal/bot/command/admin"
	"unibot/internal/bot/command/general"

	"github.com/disgoorg/disgo/events"
)

// HandlerEntry の型定義を disgo 用に更新
type HandlerEntry struct {
	// Session は不要になり、InteractionCreate は events.ApplicationCommandInteractionCreate に変更
	Handler   func(*internal.BotContext, *events.ApplicationCommandInteractionCreate)
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
}
