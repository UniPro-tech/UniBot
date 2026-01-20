package command

import "github.com/bwmarrin/discordgo"

var Handlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
	"ping": Ping,
}
