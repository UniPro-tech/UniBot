package general

import (
	"time"
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

	repo := repository.NewPinSettingRepository(ctx.DB)
	setting := &model.PinSetting{
		ID:        i.ChannelID,
		URL:       sentMessage.ID,
		Title:     "Pinned Message",
		Content:   message,
		GuildID:   i.GuildID,
		ChannelID: i.ChannelID,
	}

	if err := repo.Update(setting); err != nil {
		if err := repo.Create(setting); err != nil {
			replyPinError(s, i, config, "エラー", "ピン留めの保存に失敗しました。")
			return true
		}
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "メッセージをピン留めしました: `" + message + "`",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	return true
}

func getPinModalValue(data discordgo.ModalSubmitInteractionData, customID string) string {
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

func replyPinSuccess(s *discordgo.Session, i *discordgo.InteractionCreate, config *internal.Config, title string) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:     title,
					Color:     config.Colors.Success,
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}
