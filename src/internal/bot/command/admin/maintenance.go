package admin

import (
	"unibot/internal/bot/command/admin/maintenance"

	"github.com/disgoorg/disgo/discord"
)

func LoadMaintenanceCommandContext() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "maintenance",
		Description: "メンテナンス用コマンド",
		Options: []discord.ApplicationCommandOption{
			maintenance.LoadStatusCommandContext(),
			maintenance.LoadShutdownCommandContext(),
		},
	}
}
