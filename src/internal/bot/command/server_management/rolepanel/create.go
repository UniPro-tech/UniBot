package rolepanel

import (
	"fmt"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func LoadCreateCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "create",
		Description: "ロールパネルを作成します",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "title",
				Description: "パネルのタイトル",
				Required:    true,
				MaxLength:   256,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "description",
				Description: "パネルの説明",
				Required:    false,
				MaxLength:   4096,
			},
		},
	}
}

func Create(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	options := i.ApplicationCommandData().Options[0].Options

	var title, description string
	for _, opt := range options {
		switch opt.Name {
		case "title":
			title = opt.StringValue()
		case "description":
			description = opt.StringValue()
		}
	}

	// パネルメッセージを送信
	panelEmbed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       config.Colors.Primary,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "ロールを選択してください",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	msg, err := s.ChannelMessageSendEmbed(i.ChannelID, panelEmbed)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルの作成中にエラーが発生しました。",
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

	// データベースに保存
	repo := repository.NewRolePanelRepository(ctx.DB)
	panel := &model.RolePanel{
		GuildID:     i.GuildID,
		ChannelID:   i.ChannelID,
		MessageID:   msg.ID,
		Title:       title,
		Description: description,
	}

	err = repo.Create(panel)
	if err != nil {
		// メッセージを削除
		_ = s.ChannelMessageDelete(i.ChannelID, msg.ID)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルの保存中にエラーが発生しました。",
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

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "ロールパネルを作成しました",
					Description: fmt.Sprintf("メッセージID: `%s`\n\n`/rolepanel add` コマンドでロールを追加してください。", msg.ID),
					Color:       config.Colors.Success,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "タイトル",
							Value:  title,
							Inline: true,
						},
					},
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
