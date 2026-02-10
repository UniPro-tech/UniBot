package general

import (
	"time"
	"unibot/internal"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func LoadUnpinCommandContext() *discordgo.ApplicationCommand {
	perm := int64(discordgo.PermissionManageMessages)
	dm := false
	contexts := []discordgo.InteractionContextType{discordgo.InteractionContextGuild}
	return &discordgo.ApplicationCommand{
		Name:                     "unpin",
		Description:              "ピン留めを解除します。",
		DefaultMemberPermissions: &perm,
		DMPermission:             &dm,
		Contexts:                 &contexts,
	}
}

func Unpin(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config

	if !hasPinPermission(s, i) {
		replyPinError(s, i, config, "権限がありません", "この操作を実行する権限がありません。")
		return
	}

	repo := repository.NewPinSettingRepository(ctx.DB)
	settings, err := repo.GetByChannelID(i.ChannelID)
	if err != nil {
		replyPinError(s, i, config, "エラーが発生しました", "ピン留めの解除中にエラーが発生しました。")
		return
	}
	if len(settings) == 0 {
		content := "このチャンネルにはピン留めされたメッセージがありません。"
		_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return
	}

	err = repo.DeleteByChannelID(i.ChannelID)
	if err != nil {
		replyPinError(s, i, config, "エラーが発生しました", "ピン留めの解除中にエラーが発生しました。")
		return
	}

	successEmbed := &discordgo.MessageEmbed{
		Title:     "ピン留めを解除しました",
		Color:     config.Colors.Success,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{successEmbed},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
}
