package handler

import (
	"strconv"
	"strings"
	"time"
	"unibot/internal"
	"unibot/internal/bot/command"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(ctx *internal.BotContext) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			handleApplicationCommand(ctx, s, i)
		case discordgo.InteractionMessageComponent:
			handleMessageComponent(ctx, s, i)
		}
	}
}

func handleApplicationCommand(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name

	if h, ok := command.Handlers[name]; ok {
		h(ctx, s, i)
	}
}

func handleMessageComponent(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	config := ctx.Config

	// 辞書削除のセレクトメニュー
	if strings.HasPrefix(customID, "tts_dict_remove") {
		values := i.MessageComponentData().Values
		if len(values) == 0 {
			return
		}

		selectedID, err := strconv.ParseUint(values[0], 10, 64)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "エラー",
							Description: "不正なIDです。",
							Color:       config.Colors.Error,
							Timestamp:   time.Now().Format(time.RFC3339),
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		repo := repository.NewTTSDictionaryRepository(ctx.DB)

		// 削除対象の単語を取得
		entry, err := repo.GetByID(uint(selectedID))
		if err != nil || entry == nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "エラー",
							Description: "単語が見つかりませんでした。",
							Color:       config.Colors.Error,
							Timestamp:   time.Now().Format(time.RFC3339),
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		// セキュリティチェック: エントリが現在のギルドに属しているか確認
		if entry.GuildID != i.GuildID {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "エラー",
							Description: "この単語を削除する権限がありません。",
							Color:       config.Colors.Error,
							Timestamp:   time.Now().Format(time.RFC3339),
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		word := entry.Word

		// 削除実行
		err = repo.DeleteByID(uint(selectedID))
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "エラー",
							Description: "単語の削除に失敗しました。",
							Color:       config.Colors.Error,
							Timestamp:   time.Now().Format(time.RFC3339),
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "単語を削除しました",
						Description: "「" + word + "」を辞書から削除しました。",
						Color:       config.Colors.Success,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
	}
}
