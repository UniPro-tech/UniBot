package schedule

import (
	"fmt"
	"time"
	"unibot/internal"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func LoadListCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "list",
		Description: "予約投稿の一覧を表示します",
	}
}

func List(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	repo := repository.NewScheduleSettingRepository(ctx.DB)

	settings, err := repo.GetByChannelID(i.ChannelID)
	if err != nil {
		_ = RespondEdit(s, i, &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "スケジュールの取得に失敗しました。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	if len(settings) == 0 {
		_ = RespondEdit(s, i, &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "予約投稿一覧",
					Description: "予約投稿はまだありません。",
					Color:       config.Colors.Success,
					Timestamp:   time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:     "予約投稿一覧",
		Color:     config.Colors.Success,
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Requested by " + i.Member.DisplayName(),
			IconURL: i.Member.AvatarURL(""),
		},
	}

	for _, setting := range settings {
		repeatText := "いいえ"
		if setting.Cron != "" {
			repeatText = describeCron(setting.Cron)
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("ジョブID: %s", setting.ID),
			Value: fmt.Sprintf("メッセージ: %s\n次回実行予定: <t:%d:F>\n繰り返し: %s", setting.Content, setting.NextRunAt, repeatText),
		})
	}

	_ = RespondEdit(s, i, &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{embed},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
}
