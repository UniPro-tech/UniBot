package rolepanel

import (
	"fmt"
	"strings"
	"time"
	"unibot/internal"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func LoadListCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "list",
		Description: "このサーバーのロールパネル一覧を表示します",
	}
}

func List(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	repo := repository.NewRolePanelRepository(ctx.DB)

	panels, err := repo.ListByGuild(i.GuildID)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルの取得中にエラーが発生しました。",
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

	if len(panels) == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "ロールパネル一覧",
						Description: "このサーバーにはロールパネルがありません。\n\n`/rolepanel create` でパネルを作成してください。",
						Color:       config.Colors.Primary,
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

	var fields []*discordgo.MessageEmbedField
	for _, panel := range panels {
		var roles []string
		for _, opt := range panel.Options {
			roles = append(roles, fmt.Sprintf("<@&%s>", opt.RoleID))
		}

		roleList := "なし"
		if len(roles) > 0 {
			roleList = strings.Join(roles, ", ")
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s (ID: %s)", panel.Title, panel.MessageID),
			Value:  fmt.Sprintf("チャンネル: <#%s>\nロール: %s", panel.ChannelID, roleList),
			Inline: false,
		})
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "ロールパネル一覧",
					Description: fmt.Sprintf("このサーバーには %d 個のロールパネルがあります。", len(panels)),
					Color:       config.Colors.Primary,
					Fields:      fields,
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
