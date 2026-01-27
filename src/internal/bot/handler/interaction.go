package handler

import (
	"strings"
	"unibot/internal"
	"unibot/internal/bot/command"
	"unibot/internal/bot/command/general/schedule"
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
		case discordgo.InteractionModalSubmit:
			handleModalSubmit(ctx, s, i)
		}
	}
}

func handleApplicationCommand(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if h, ok := command.Handlers[name]; ok {
		h(ctx, s, i)
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

func handleModalSubmit(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if schedule.HandleModalSubmit(ctx, s, i) {
		return
	}
}

