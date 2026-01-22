package handler

import (
	"unibot/internal"
	"unibot/internal/bot/command"

	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(ctx *internal.BotContext) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		name := i.ApplicationCommandData().Name

		if h, ok := command.Handlers[name]; ok {
			h(ctx, s, i)
		}
	}
}
