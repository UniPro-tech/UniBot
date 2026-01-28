package command

import (
	"unibot/internal/bot/command/admin"
	"unibot/internal/bot/command/general"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	general.LoadPingCommandContext(),
	general.LoadAboutCommandContext(),
	general.LoadPinCommandContext(),
	general.LoadUnpinCommandContext(),
	general.LoadPinSelectCommandContext(),
	general.LoadTtsCommandContext(),
	general.LoadHelpCommandContext(),
	general.LoadScheduleCommandContext(),
	admin.LoadMaintenanceCommandContext(),
}
