package dict

import (
	"fmt"
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func LoadRemoveCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "remove",
		Description: "TTS辞書から単語を削除します",
	}
}

func Remove(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	repo := repository.NewTTSDictionaryRepository(ctx.DB)

	done := make(chan struct{})
	go func() {
		select {
		case <-done:
			return
		case <-time.After(3 * time.Minute):
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "ボイスチャンネルの情報を取得できませんでした。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			})
			if err != nil {
				log.Println("Failed to edit deferred interaction on timeout:", err)
			}
		}
	}()
	defer close(done)

	// 管理者かどうか確認
	perms, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
						Title:       "エラー",
						Description: "権限の確認中にエラーが発生しました。",
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

	isAdmin := perms&discordgo.PermissionAdministrator != 0

	// 辞書エントリを取得
	var entries []*model.TTSDictionary
	if isAdmin {
		entries, err = repo.ListByGuild(i.GuildID)
	} else {
		entries, err = repo.ListByGuildUser(i.GuildID, i.Member.User.ID)
	}

	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
						Title:       "エラー",
						Description: "辞書の取得中にエラーが発生しました。",
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

	if len(entries) == 0 {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
						Title:       "辞書が空です",
						Description: "辞書に登録されている単語がありません。",
						Color:       config.Colors.Warning,
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

	// 25件以上ある場合は最初の25件のみ表示
	displayEntries := entries
	if len(entries) > 25 {
		displayEntries = entries[:25]
	}

	// セレクトメニューを作成
	options := make([]discordgo.SelectMenuOption, len(displayEntries))
	for idx, entry := range displayEntries {
		options[idx] = discordgo.SelectMenuOption{
			Label:       entry.Word,
			Value:       fmt.Sprintf("%d", entry.ID),
			Description: entry.Definition,
		}
	}

	content := "削除したい単語を選んでください。"
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "tts_dict_remove",
					Placeholder: "削除する単語を選んでください",
					Options:     options,
				},
			},
		},
	}

	// InteractionResponseEdit を実行
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:    &content,
		Components: &components,
	})
}
