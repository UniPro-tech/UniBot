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
	} else if isTtsSetVoice(i) {
		response.Data = &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		}
	}
	s.InteractionRespond(i.Interaction, response)
	if entry, ok := command.Handlers[name]; ok {
		entry.Handler(ctx, s, i)
	}
}

func isTtsSetVoice(i *discordgo.InteractionCreate) bool {
	if i.ApplicationCommandData().Name != "tts" {
		return false
	}
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		return false
	}
	group := options[0]
	if group.Type != discordgo.ApplicationCommandOptionSubCommandGroup || group.Name != "set" {
		return false
	}
	if len(group.Options) == 0 {
		return false
	}
	sub := group.Options[0]
	return sub.Type == discordgo.ApplicationCommandOptionSubCommand && sub.Name == "voice"
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
