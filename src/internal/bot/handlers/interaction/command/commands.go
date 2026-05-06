package command

import (
	"unibot/internal/bot/handlers/interaction/command/admin"
	"unibot/internal/bot/handlers/interaction/command/general"

	"github.com/disgoorg/disgo/discord"
)

var GeneralCommands = []discord.ApplicationCommandCreate{
	general.LoadPingCommandContext(),
	general.LoadAboutCommandContext(),
	general.LoadColorCodeCommandContext(),
	general.LoadHelpCommandContext(),
	general.LoadTtsCommandContext(),
}

var AdminCommands = []discord.ApplicationCommandCreate{
	admin.LoadMaintenanceCommandContext(),
}
