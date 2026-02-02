package handler

import (
	"strings"
	"unibot/internal"
	"unibot/internal/bot/command"
	"unibot/internal/bot/messageComponent"

	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(ctx *internal.BotContext) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			handleApplicationCommand(ctx, s, i)
		case discordgo.InteractionMessageComponent:
			handleMessageComponent(ctx, s, i)
		}
	}
}

func handleApplicationCommand(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}
	if entry, ok := command.Handlers[name]; ok && entry.Ephemeral {
		response.Data = &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		}
	}
	s.InteractionRespond(i.Interaction, response)
	if entry, ok := command.Handlers[name]; ok {
		entry.Handler(ctx, s, i)
	}
}

func handleMessageComponent(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID

	for prefix, handler := range messageComponent.Handlers {
		if strings.HasPrefix(customID, prefix) {
			handler(ctx, s, i)
			return
		}
	}
}
