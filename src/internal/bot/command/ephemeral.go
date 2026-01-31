package command

import "github.com/bwmarrin/discordgo"

// ShouldDeferEphemeral returns true when the initial deferred response should be ephemeral.
func ShouldDeferEphemeral(i *discordgo.InteractionCreate) bool {
	switch i.ApplicationCommandData().Name {
	case "maintenance":
		return true
	default:
		return false
	}
}
