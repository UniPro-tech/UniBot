package general

import (
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

// HandlePinModalSubmit はピン留めモーダルの送信を処理する
func HandlePinModalSubmit(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	data := i.ModalSubmitData()
	if data.CustomID != "pin_message" {
		return false
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	config := ctx.Config
	if !hasPinPermission(s, i) {
		replyPinError(s, i, config, "権限がありません", "この操作を実行する権限がありません。")
		return true
	}

	message := getPinModalValue(data, "message")
	if message == "" {
		replyPinError(s, i, config, "入力エラー", "投稿内容を入力してください。")
		return true
	}

	channel, err := s.State.Channel(i.ChannelID)
	if err != nil {
		channel, _ = s.Channel(i.ChannelID)
	}
	if channel == nil || channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM {
		replyPinError(s, i, config, "エラー", "このチャンネルではメッセージを送信できません。")
		return true
	}

	repo := repository.NewPinSettingRepository(ctx.DB)
	existing, err := repo.GetByChannelID(i.ChannelID)
	if err != nil {
		replyPinError(s, i, config, "エラー", "ピン留めの取得に失敗しました。")
		return true
	}
	if len(existing) > 0 {
		replyPinError(s, i, config, "エラー", "このチャンネルには既にピン留めされたメッセージがあります。\n最初にそれを`/unpin`で解除してください。")
		return true
	}

	embed := &discordgo.MessageEmbed{
		Description: message,
		Color:       config.Colors.Success,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Pinned Message",
		},
	}

	sentMessage, err := s.ChannelMessageSendEmbed(i.ChannelID, embed)
	if err != nil {
		replyPinError(s, i, config, "エラー", "メッセージの送信に失敗しました。")
		return true
	}

	setting := &model.PinSetting{
		ID:        i.ChannelID,
		URL:       sentMessage.ID,
		Title:     "Pinned Message",
		Content:   message,
		GuildID:   i.GuildID,
		ChannelID: i.ChannelID,
	}

	if err := repo.Create(setting); err != nil {
		replyPinError(s, i, config, "エラー", "ピン留めの保存に失敗しました。")
		return true
	}

	content := "メッセージをピン留めしました: `" + message + "`"
	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})

	return true
}

func getPinModalValue(data discordgo.ModalSubmitInteractionData, customID string) string {
	for _, comp := range data.Components {
		switch row := comp.(type) {
		case *discordgo.ActionsRow:
			if value := getTextInputValue(row.Components, customID); value != "" {
				return value
			}
		case discordgo.ActionsRow:
			if value := getTextInputValue(row.Components, customID); value != "" {
				return value
			}
		}
	}

	return ""
}

func getTextInputValue(components []discordgo.MessageComponent, customID string) string {
	for _, component := range components {
		switch input := component.(type) {
		case *discordgo.TextInput:
			if input.CustomID == customID {
				return input.Value
			}
		case discordgo.TextInput:
			if input.CustomID == customID {
				return input.Value
			}
		}
	}
	return ""
}
