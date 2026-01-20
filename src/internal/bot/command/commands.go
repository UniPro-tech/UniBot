package command

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "スピードテストを行います",
	},
	{
		Name:        "about",
		Description: "ボットの情報を表示します",
	},
}
