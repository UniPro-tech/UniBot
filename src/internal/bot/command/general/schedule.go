package general

import (
	"time"
	"unibot/internal"
	schedulecmd "unibot/internal/bot/command/general/schedule"

	"github.com/bwmarrin/discordgo"
)

func LoadScheduleCommandContext() *discordgo.ApplicationCommand {
	perm := int64(discordgo.PermissionManageMessages)
	return &discordgo.ApplicationCommand{
		Name:                     "schedule",
		Description:              "予約投稿を管理します",
		DefaultMemberPermissions: &perm,
		Options: []*discordgo.ApplicationCommandOption{
			schedulecmd.LoadSetCommandContext(),
			schedulecmd.LoadListCommandContext(),
			schedulecmd.LoadRemoveCommandContext(),
		},
	}
}

var scheduleHandler = map[string]func(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate){
	"set":    schedulecmd.Set,
	"list":   schedulecmd.List,
	"remove": schedulecmd.Remove,
}

func Schedule(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config

	if i.GuildID == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このコマンドはサーバー内でのみ使用できます。",
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
		return
	}

	perms, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil || perms&discordgo.PermissionManageMessages == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このコマンドを実行する権限がありません。",
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
		return
	}

	subCommand := i.ApplicationCommandData().Options[0]
	if handler, exists := scheduleHandler[subCommand.Name]; exists {
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
