package general

import (
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func LoadPinSelectCommandContext() *discordgo.ApplicationCommand {
	contexts := []discordgo.InteractionContextType{discordgo.InteractionContextGuild}
	return &discordgo.ApplicationCommand{
		Name:     "Pinするメッセージを選択",
		Type:     discordgo.MessageApplicationCommand,
		Contexts: &contexts,
	}
}

func PinSelect(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config

	if !hasPinPermission(s, i) {
		replyPinError(s, i, config, "権限がありません", "この操作を実行する権限がありません。")
		return
	}

	data := i.ApplicationCommandData()
	if data.Resolved == nil || data.Resolved.Messages == nil {
		replyPinError(s, i, config, "エラー", "メッセージの取得に失敗しました。")
		return
	}

	targetMsg, ok := data.Resolved.Messages[data.TargetID]
	if !ok || targetMsg == nil {
		replyPinError(s, i, config, "エラー", "メッセージの取得に失敗しました。")
		return
	}

	if targetMsg.Author != nil && targetMsg.Author.Bot {
		replyPinError(s, i, config, "エラー", "ボットのメッセージはピン留めできません。")
		return
	}

	channel, err := s.State.Channel(i.ChannelID)
	if err != nil {
		channel, _ = s.Channel(i.ChannelID)
	}
	if channel == nil || channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM {
		replyPinError(s, i, config, "エラー", "このチャンネルではメッセージをピン留めできません。")
		return
	}

	repo := repository.NewPinSettingRepository(ctx.DB)
	settings, err := repo.GetByChannelID(i.ChannelID)
	if err != nil {
		replyPinError(s, i, config, "エラー", "ピン留めの取得に失敗しました。")
		return
	}
	if len(settings) > 0 {
		replyPinError(s, i, config, "エラー", "このチャンネルには既にピン留めされたメッセージがあります。\n最初にそれを`/unpin`で解除してください。")
		return
	}

	embed := &discordgo.MessageEmbed{
		Description: targetMsg.Content,
		Color:       config.Colors.Success,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Pinned Message",
		},
	}

	sentMessage, err := s.ChannelMessageSendEmbed(i.ChannelID, embed)
	if err != nil {
		replyPinError(s, i, config, "エラー", "メッセージの送信に失敗しました。")
		return
	}

	setting := &model.PinSetting{
		ID:        i.ChannelID,
		URL:       sentMessage.ID,
		Title:     "Pinned Message",
		Content:   targetMsg.Content,
		GuildID:   i.GuildID,
		ChannelID: i.ChannelID,
	}

	err = repo.Create(setting)
	if err != nil {
		replyPinError(s, i, config, "エラー", "ピン留めの保存に失敗しました。")
		return
	}

	successEmbed := &discordgo.MessageEmbed{
		Title:       "メッセージをピン留めしました",
		Description: "このメッセージは今後ピン留めされます。\nファイルは保存されないのでご注意ください。",
		Color:       config.Colors.Success,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{successEmbed},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
}
