package general

import (
	"unibot/internal/bot/handlers/interaction/command/general/rss"

	"github.com/disgoorg/disgo/discord"
)

func LoadRssCommandContext() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "rss",
		Description: "RSSフィードを受信します",
		Options: []discord.ApplicationCommandOption{
			rss.LoadSubscribeCommandContext(),
		},
	}
}
