package command

import (
	"unibot/internal/bot/command/admin"
	"unibot/internal/bot/command/general"
	"unibot/internal/bot/command/server_management"

	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	general.LoadPingCommandContext(),
	general.LoadAboutCommandContext(),
	general.LoadColorCodeCommandContext(),
	general.LoadTtsCommandContext(),
	general.LoadHelpCommandContext(),
	admin.LoadMaintenanceCommandContext(),
	server_management.LoadRolepanelCommandContext(),
}
