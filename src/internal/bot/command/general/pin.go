package general

import (
	"time"
	"unibot/internal"

	"github.com/bwmarrin/discordgo"
)

func LoadPinCommandContext() *discordgo.ApplicationCommand {
	perm := int64(discordgo.PermissionManageMessages)
	dm := false
	contexts := []discordgo.InteractionContextType{discordgo.InteractionContextGuild}
	return &discordgo.ApplicationCommand{
		Name:                     "pin",
		Description:              "メッセージをピン留めします。",
		DefaultMemberPermissions: &perm,
		DMPermission:             &dm,
		Contexts:                 &contexts,
	}
}

func Pin(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config

	if !hasPinPermission(s, i) {
		replyPinError(s, i, config, "権限がありません", "この操作を実行する権限がありません。")
		return
	}

	showPinModal(s, i)
}

func showPinModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "pin_message",
			Title:    "メッセージのピン留め",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "message",
						Label:       "投稿内容",
						Style:       discordgo.TextInputParagraph,
						Placeholder: "投稿内容を入力してください。すでにPinされたメッセージがある場合は上書きされます。",
						Required:    true,
					},
				}},
			},
		},
	})
}

func hasPinPermission(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if i.Member == nil || i.Member.User == nil {
		return false
	}
	perms, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		return false
	}
	return perms&discordgo.PermissionManageMessages != 0
}

func replyPinError(s *discordgo.Session, i *discordgo.InteractionCreate, config *internal.Config, title, description string) {
	footer := &discordgo.MessageEmbedFooter{Text: "Requested by Unknown"}
	if i.Member != nil {
		footer.Text = "Requested by " + i.Member.DisplayName()
		footer.IconURL = i.Member.AvatarURL("")
	} else if i.User != nil {
		footer.Text = "Requested by " + i.User.Username
		footer.IconURL = i.User.AvatarURL("")
	}

	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       title,
				Description: description,
				Color:       config.Colors.Error,
				Footer:      footer,
				Timestamp: time.Now().Format(time.RFC3339),
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
}
