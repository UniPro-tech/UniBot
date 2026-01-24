package rolepanel

import (
	"fmt"
	"time"
	"unibot/internal"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func LoadAddCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "add",
		Description: "ロールパネルにロールを追加します",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "追加するロール",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "label",
				Description: "セレクトメニューに表示するラベル",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "description",
				Description: "ロールの説明",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "emoji",
				Description: "絵文字 (例: 🎮)",
				Required:    false,
			},
		},
	}
}

func Add(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	options := i.ApplicationCommandData().Options[0].Options

	var roleID, label, description, emoji string
	for _, opt := range options {
		switch opt.Name {
		case "role":
			roleID = opt.RoleValue(s, i.GuildID).ID
		case "label":
			label = opt.StringValue()
		case "description":
			description = opt.StringValue()
		case "emoji":
			emoji = opt.StringValue()
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

	if len(panels) == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このサーバーにはロールパネルがありません。\n先に `/rolepanel create` でパネルを作成してください。",
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
	for _, panel := range panels {
		selectOptions = append(selectOptions, discordgo.SelectMenuOption{
			Label:       panel.Title,
			Value:       panel.MessageID,
			Description: fmt.Sprintf("%d個のロール", len(panel.Options)),
		})
	}

	// CustomIDにロール情報をエンコード (roleID|label|description|emoji)
	customID := fmt.Sprintf("rolepanel_add_%s|%s|%s|%s", roleID, label, description, emoji)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "パネルを選択",
					Description: fmt.Sprintf("ロール <@&%s> を追加するパネルを選択してください。", roleID),
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
