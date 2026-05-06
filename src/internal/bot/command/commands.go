package command

import (
	"unibot/internal/bot/command/admin"
	"unibot/internal/bot/command/general"

	"github.com/disgoorg/disgo/discord"
)

var GeneralCommands = []discord.ApplicationCommandCreate{
	general.LoadPingCommandContext(),
	general.LoadAboutCommandContext(),
	general.LoadColorCodeCommandContext(),
	general.LoadHelpCommandContext(),
	//general.LoadTtsCommandContext(),
	//general.LoadHelpCommandContext(),
}

var AdminCommands = []discord.ApplicationCommandCreate{
	admin.LoadMaintenanceCommandContext(),
}
