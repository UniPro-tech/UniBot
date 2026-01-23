package messageComponent

import (
	"strconv"
	"time"
	"unibot/internal"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterHandler("tts_dict_remove", HandleTTSDictRemove)
}

// HandleTTSDictRemove は辞書削除のセレクトメニューを処理します
func HandleTTSDictRemove(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
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
