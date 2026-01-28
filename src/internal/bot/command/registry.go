package command

import (
	"unibot/internal"
	"unibot/internal/bot/command/admin"
	"unibot/internal/bot/command/general"

	"github.com/bwmarrin/discordgo"
)

var Handlers = map[string]func(*internal.BotContext, *discordgo.Session, *discordgo.InteractionCreate){
	"ping":        general.Ping,
	"about":       general.About,
	"maintenance": admin.Maintenance,
	"tts":         general.Tts,
	"help":		   general.Help,
}
