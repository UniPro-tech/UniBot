package messageComponent

import (
	"log"
	"strconv"
	"strings"
	"time"
	"unibot/internal"
	"unibot/internal/bot/ttsutil"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterHandler(ttsutil.VoiceSelectCustomID, HandleTTSSetVoice)
	RegisterHandler(ttsutil.VoicePageCustomIDPrefix, HandleTTSSetVoicePage)
}

// HandleTTSSetVoice は話者選択のセレクトメニューを処理します
func HandleTTSSetVoice(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	values := i.MessageComponentData().Values
	if len(values) == 0 {
		return
	}

	speakerID := values[0]
	memberID, _, _ := ttsutil.GetInteractionUser(i)
	if memberID == "" {
		log.Println("HandleTTSSetVoice: missing user information on interaction")
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "ユーザー情報の取得に失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		}); err != nil {
			log.Println("Failed to respond interaction:", err)
		}
		return
	}

	memberRepo := repository.NewMemberRepository(ctx.DB)
	if err := memberRepo.Create(memberID); err != nil {
		log.Println("Error creating member:", err)
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "メンバー情報の作成に失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		}); err != nil {
			log.Println("Failed to respond interaction:", err)
		}
		return
	}

	repo := repository.NewTTSPersonalSettingRepository(ctx.DB)
	setting, err := repo.GetByMember(memberID)
	if err != nil {
		log.Println("Error fetching TTS personal setting:", err)
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "TTS個人設定の取得に失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		}); err != nil {
			log.Println("Failed to respond interaction:", err)
		}
		return
	}

	if setting == nil {
		defaultSetting := repository.DefaultTTSPersonalSetting
		setting = &defaultSetting
		setting.MemberID = memberID
		setting.SpeakerID = speakerID
		if err := repo.Create(setting); err != nil {
			log.Println("Error creating TTS personal setting:", err)
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "エラー",
							Description: "TTS個人設定の作成に失敗しました。",
							Color:       config.Colors.Error,
							Timestamp:   time.Now().Format(time.RFC3339),
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			}); err != nil {
				log.Println("Failed to respond interaction:", err)
			}
			return
		}
	} else {
		setting.SpeakerID = speakerID
		if err := repo.Update(setting); err != nil {
			log.Println("Error updating TTS personal setting:", err)
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "エラー",
							Description: "TTS個人設定の更新に失敗しました。",
							Color:       config.Colors.Error,
							Timestamp:   time.Now().Format(time.RFC3339),
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			}); err != nil {
				log.Println("Failed to respond interaction:", err)
			}
			return
		}
	}

	label := ttsutil.ResolveSpeakerLabel(ctx, speakerID)
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "話者設定を更新しました",
					Description: "選択した話者: " + label,
					Color:       config.Colors.Success,
					Timestamp:   time.Now().Format(time.RFC3339),
				},
			},
			Components: []discordgo.MessageComponent{},
		},
	}); err != nil {
		log.Println("Failed to respond interaction:", err)
	}
}

// HandleTTSSetVoicePage は話者選択のページ送りを処理します
func HandleTTSSetVoicePage(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	parts := strings.SplitN(customID, ":", 2)
	if len(parts) != 2 {
		return
	}

	pageIndex, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	speakers, err := ttsutil.FetchSpeakers(ctx)
	if err != nil {
		log.Println("Failed to fetch speakers:", err)
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "話者情報の取得に失敗しました。",
						Color:       ctx.Config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		}); err != nil {
			log.Println("Failed to respond interaction:", err)
		}
		return
	}

	pages := ttsutil.BuildSpeakerPages(speakers, ttsutil.SpeakerPageSize)
	if len(pages) == 0 {
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "話者情報が取得できませんでした。",
						Color:       ctx.Config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		}); err != nil {
			log.Println("Failed to respond interaction:", err)
		}
		return
	}

	memberID, _, _ := ttsutil.GetInteractionUser(i)
	if memberID == "" {
		log.Println("HandleTTSSetVoicePage: missing user information on interaction")
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "ユーザー情報の取得に失敗しました。",
						Color:       ctx.Config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		}); err != nil {
			log.Println("Failed to respond interaction:", err)
		}
		return
	}
	currentSpeakerID := ttsutil.GetCurrentSpeakerID(ctx, memberID)
	content, components := ttsutil.BuildVoiceMessage(pageIndex, pages, currentSpeakerID)

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Components: components,
		},
	}); err != nil {
		log.Println("Failed to respond interaction:", err)
	}
}
