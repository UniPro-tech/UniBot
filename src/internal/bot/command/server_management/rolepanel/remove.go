package rolepanel

import (
	"fmt"
	"time"
	"unibot/internal"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func LoadRemoveCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "remove",
		Description: "ロールパネルからロールを削除します",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "削除するロール",
				Required:    true,
			},
		},
	}
}

func Remove(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	options := i.ApplicationCommandData().Options[0].Options

	var roleID string
	for _, opt := range options {
		if opt.Name == "role" {
			roleID = opt.RoleValue(s, i.GuildID).ID
		}
	}

	repo := repository.NewRolePanelRepository(ctx.DB)

	// このギルドのパネル一覧を取得
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

	// このロールを含むパネルをフィルタ
	var matchingPanels []*struct {
		Title     string
		MessageID string
		RoleCount int
	}
	hadReadError := false
	for _, panel := range panels {
		roleIDsByKey, err := loadPanelRoleIDs(s, panel)
		if err != nil {
			hadReadError = true
			continue
		}

		for _, currentRoleID := range roleIDsByKey {
			if currentRoleID != roleID {
				continue
			}

			matchingPanels = append(matchingPanels, &struct {
				Title     string
				MessageID string
				RoleCount int
			}{
				Title:     panel.Title,
				MessageID: panel.MessageID,
				RoleCount: len(panel.Options),
			})
			break
		}
	}

	if len(matchingPanels) == 0 {
		description := fmt.Sprintf("ロール <@&%s> を含むパネルが見つかりません。", roleID)
		if hadReadError {
			description += "\n一部のパネルを確認できませんでした。"
		}
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: description,
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

	// パネル選択用セレクトメニューを作成
	var selectOptions []discordgo.SelectMenuOption
	for _, panel := range matchingPanels {
		selectOptions = append(selectOptions, discordgo.SelectMenuOption{
			Label:       panel.Title,
			Value:       panel.MessageID,
			Description: fmt.Sprintf("%d個のロール", panel.RoleCount),
		})
	}
	customID := fmt.Sprintf("rolepanel_remove_%s", roleID)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "パネルを選択",
					Description: fmt.Sprintf("ロール <@&%s> を削除するパネルを選択してください。", roleID),
					Color:       config.Colors.Primary,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    customID,
							Placeholder: "パネルを選択...",
							Options:     selectOptions,
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}
