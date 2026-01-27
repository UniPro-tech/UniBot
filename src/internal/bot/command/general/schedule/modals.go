package schedule

import (
	"fmt"
	"strings"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"
	"unibot/internal/scheduler"

	"github.com/bwmarrin/discordgo"
)

// モーダル送信を処理する
func HandleModalSubmit(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	data := i.ModalSubmitData()

	switch data.CustomID {
	case "schedule_create_onetime":
		handleCreateOnetime(ctx, s, i, data)
		return true
	case "schedule_create_repeat":
		handleCreateRepeat(ctx, s, i, data)
		return true
	default:
		return false
	}
}

func handleCreateOnetime(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ModalSubmitInteractionData) {
	config := ctx.Config

	if !hasManageMessagesPermission(s, i) {
		replyPermissionError(s, i, config)
		return
	}

	message := getTextInputValue(data, "message")
	timeText := getTextInputValue(data, "time")

	if message == "" || timeText == "" {
		replyError(s, i, config, "入力エラー", "投稿内容と投稿時間を入力してください。")
		return
	}

	jst := scheduler.JST()
	scheduledTime, err := time.ParseInLocation("2006-01-02 15:04", strings.TrimSpace(timeText), jst)
	if err != nil {
		replyError(s, i, config, "時間の形式が正しくありません。", "YYYY-MM-DD HH:mm の形式で入力してください。")
		return
	}

	if scheduledTime.Before(time.Now().In(jst)) {
		replyError(s, i, config, "過去の日時は指定できません。", "未来の日時を入力してください。")
		return
	}

	setting := &model.ScheduleSetting{
		ID:        i.ID,
		ChannelID: i.ChannelID,
		Content:   message,
		NextRunAt: scheduledTime.Unix(),
		Cron:      "",
		GuildID:   i.GuildID,
		AuthorID:  i.Member.User.ID,
	}

	repo := repository.NewScheduleSettingRepository(ctx.DB)
	if err := repo.Create(setting); err != nil {
		replyError(s, i, config, "エラー", "スケジュールの作成に失敗しました。")
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("メッセージを<t:%d:F>に送信するようにスケジュールしました。(ジョブID: %s)", scheduledTime.Unix(), i.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func handleCreateRepeat(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ModalSubmitInteractionData) {
	config := ctx.Config

	if !hasManageMessagesPermission(s, i) {
		replyPermissionError(s, i, config)
		return
	}

	message := getTextInputValue(data, "message")
	inputText := getTextInputValue(data, "time")

	if message == "" || inputText == "" {
		replyError(s, i, config, "入力エラー", "投稿内容と時間を入力してください。")
		return
	}

	jst := scheduler.JST()
	cronText, err := convertToCron(strings.TrimSpace(inputText))
	if err != nil {
		replyError(s, i, config, "時間の形式が不正です。", "時間の形式が正しくありません。もう一度確認してください。")
		return
	}

	nextRunAt, err := scheduler.NextRunAtFromCron(cronText, time.Now().In(jst))
	if err != nil {
		replyError(s, i, config, "時間の形式が不正です。", "時間の形式が正しくありません。もう一度確認してください。")
		return
	}

	setting := &model.ScheduleSetting{
		ID:        i.ID,
		ChannelID: i.ChannelID,
		Content:   message,
		NextRunAt: nextRunAt.Unix(),
		Cron:      strings.TrimSpace(cronText),
		GuildID:   i.GuildID,
		AuthorID:  i.Member.User.ID,
	}

	repo := repository.NewScheduleSettingRepository(ctx.DB)
	if err := repo.Create(setting); err != nil {
		replyError(s, i, config, "エラー", "スケジュールの作成に失敗しました。")
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("メッセージを%sに送信するようにスケジュールしました。(ジョブID: %s)", strings.TrimSpace(inputText), i.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func getTextInputValue(data discordgo.ModalSubmitInteractionData, customID string) string {
	for _, comp := range data.Components {
		switch row := comp.(type) {
		case *discordgo.ActionsRow:
			for _, component := range row.Components {
				if input, ok := component.(*discordgo.TextInput); ok {
					if input.CustomID == customID {
						return input.Value
					}
				}
				if input, ok := component.(discordgo.TextInput); ok {
					if input.CustomID == customID {
						return input.Value
					}
				}
			}
		case discordgo.ActionsRow:
			for _, component := range row.Components {
				if input, ok := component.(*discordgo.TextInput); ok {
					if input.CustomID == customID {
						return input.Value
					}
				}
				if input, ok := component.(discordgo.TextInput); ok {
					if input.CustomID == customID {
						return input.Value
					}
				}
			}
		}
	}

	return ""
}

func hasManageMessagesPermission(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if i.Member == nil || i.Member.User == nil {
		return false
	}
	perms, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		return false
	}
	return perms&discordgo.PermissionManageMessages != 0
}

func replyPermissionError(s *discordgo.Session, i *discordgo.InteractionCreate, config *internal.Config) {
	replyError(s, i, config, "エラー", "この操作を実行する権限がありません。")
}

func replyError(s *discordgo.Session, i *discordgo.InteractionCreate, config *internal.Config, title, description string) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       title,
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
}
