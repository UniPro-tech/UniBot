package handler

import (
	"unibot/internal/bot/command"

	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	name := i.ApplicationCommandData().Name

	if h, ok := command.Handlers[name]; ok {
		h(s, i)
	}
}
