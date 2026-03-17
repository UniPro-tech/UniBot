package server_management

import (
	"time"
	"unibot/internal"
	"unibot/internal/bot/command/server_management/rolepanel"

	"github.com/bwmarrin/discordgo"
)

func LoadRolepanelCommandContext() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:                     "rolepanel",
		Description:              "ロールパネルを管理します",
		DefaultMemberPermissions: ptrInt64(discordgo.PermissionManageRoles),
		Options: []*discordgo.ApplicationCommandOption{
			rolepanel.LoadCreateCommandContext(),
			rolepanel.LoadDeleteCommandContext(),
			rolepanel.LoadAddCommandContext(),
			rolepanel.LoadRemoveCommandContext(),
			rolepanel.LoadListCommandContext(),
		},
	}
}

func ptrInt64(i int64) *int64 {
	return &i
}

var rolepanelHandler = map[string]func(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate){
	"create": rolepanel.Create,
	"delete": rolepanel.Delete,
	"add":    rolepanel.Add,
	"remove": rolepanel.Remove,
	"list":   rolepanel.List,
}

func Rolepanel(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	subCommand := i.ApplicationCommandData().Options[0]

	if handler, exists := rolepanelHandler[subCommand.Name]; exists {
		handler(ctx, s, i)
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "不明なサブコマンドです。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}
