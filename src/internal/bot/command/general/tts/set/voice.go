package set

import (
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/bot/ttsutil"

	"github.com/bwmarrin/discordgo"
)

func LoadVoiceCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "voice",
		Description: "読み上げの話者を設定します",
	}
}

func Voice(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	memberID, requesterName, requesterAvatar := ttsutil.GetInteractionUser(i)
	if requesterName == "" {
		log.Println("Voice: missing user information on interaction")
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "ユーザー情報の取得に失敗しました。",
					Color:       config.Colors.Error,
					Timestamp:   time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			log.Println("Failed to edit deferred interaction:", err)
		}
		return
	}

	userVoiceState, err := s.State.VoiceState(i.GuildID, memberID)
	if err != nil {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "ボイスチャンネルの情報を取得できませんでした。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + requesterName,
						IconURL: requesterAvatar,
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			log.Println("Failed to edit deferred interaction:", err)
		}
		return
	}
	if userVoiceState == nil || userVoiceState.ChannelID == "" {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "先にボイスチャンネルに参加してください。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + requesterName,
						IconURL: requesterAvatar,
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			log.Println("Failed to edit deferred interaction:", err)
		}
		return
	}

	done := make(chan struct{})
	defer close(done)

	timeout := time.NewTimer(3 * time.Minute)
	defer timeout.Stop()

	go func() {
		select {
		case <-done:
			return
		case <-timeout.C:
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "処理がタイムアウトしました。もう一度お試しください。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + requesterName,
							IconURL: requesterAvatar,
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

	speakers, err := ttsutil.FetchSpeakers(ctx)
	if err != nil {
		log.Println("Failed to fetch speakers:", err)
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "話者情報の取得に失敗しました。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + requesterName,
						IconURL: requesterAvatar,
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			log.Println("Failed to edit deferred interaction:", err)
		}
		return
	}

	pages := ttsutil.BuildSpeakerPages(speakers, ttsutil.SpeakerPageSize)
	if len(pages) == 0 {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "話者情報が見つかりませんでした。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + requesterName,
						IconURL: requesterAvatar,
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			log.Println("Failed to edit deferred interaction:", err)
		}
		return
	}

	currentSpeakerID := ttsutil.GetCurrentSpeakerID(ctx, memberID)
	content, components := ttsutil.BuildVoiceMessage(0, pages, currentSpeakerID)

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:    &content,
		Components: &components,
	})
	if err != nil {
		log.Println("Failed to edit deferred interaction:", err)
	}
}
