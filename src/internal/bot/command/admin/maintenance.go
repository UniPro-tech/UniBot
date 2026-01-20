package admin

import (
	"strconv"
	"time"
	"unibot/internal"
	"unibot/internal/bot/command/admin/maintenance"

	"github.com/bwmarrin/discordgo"
)

func LoadMaintenanceCommandContext() *discordgo.ApplicationCommand {
	config := internal.LoadConfig()
	return &discordgo.ApplicationCommand{
		Name:        "maintenance",
		Description: "メンテナンス用コマンド",
		GuildID:     config.AdminGuildID,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "status",
				Description: "Botのステータスを変更",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "set",
						Description: "Botのステータスを設定します",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "text",
								Required:    true,
								Description: "ステータスメッセージ",
							},
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "status",
								Description: "ステータス種類",
								Required:    true,
								Choices: []*discordgo.ApplicationCommandOptionChoice{
									{Name: "playing", Value: strconv.Itoa(int(discordgo.ActivityTypeGame))},
									{Name: "streaming", Value: strconv.Itoa(int(discordgo.ActivityTypeStreaming))},
									{Name: "listening", Value: strconv.Itoa(int(discordgo.ActivityTypeListening))},
									{Name: "watching", Value: strconv.Itoa(int(discordgo.ActivityTypeWatching))},
									{Name: "competing", Value: strconv.Itoa(int(discordgo.ActivityTypeCompeting))},
									{Name: "custom", Value: strconv.Itoa(int(discordgo.ActivityTypeCustom))},
								},
							},
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "type",
								Description: "オンライン状態",
								Required:    true,
								Choices: []*discordgo.ApplicationCommandOptionChoice{
									{Name: "online", Value: "online"},
									{Name: "idle", Value: "idle"},
									{Name: "dnd", Value: "dnd"},
									{Name: "invisible", Value: "invisible"},
								},
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "reset",
						Description: "Botのステータスをリセットします",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "reboot",
				Description: "Botを再起動します",
			},
		},
	}
}

func IsOwner(member discordgo.Member) bool {
	config := internal.LoadConfig()
	adminRoleID := config.AdminRoleID
	for _, roleID := range member.Roles {
		if roleID == adminRoleID {
			return true
		}
	}
	return false
}

var maintenanceHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"status": maintenance.Status,
	"reboot": maintenance.Reboot,
}

func Maintenance(s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := internal.LoadConfig()
	if !IsOwner(*i.Member) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "権限エラー",
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

	subCommandGroup := i.ApplicationCommandData().Options[0]
	if subCommandGroup.Type == discordgo.ApplicationCommandOptionSubCommandGroup {
		if handler, exists := maintenanceHandler[subCommandGroup.Name]; exists {
			handler(s, i)
			return
		}
	} else {
		if handler, exists := maintenanceHandler[subCommandGroup.Name]; exists {
			handler(s, i)
			return
		}
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
