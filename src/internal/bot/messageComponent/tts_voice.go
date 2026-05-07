package messageComponent

import (
	"fmt"
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/bot/ttsutil"
	"unibot/internal/repository"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// HandleTTSSetVoice は話者選択のセレクトメニューを処理します
func HandleTTSSetVoice(ctx *internal.BotContext) func(data discord.SelectMenuInteractionData, e *handler.ComponentEvent) error {
	return func(_ discord.SelectMenuInteractionData, e *handler.ComponentEvent) error {
		config := ctx.Config
		values := e.StringSelectMenuInteractionData().Values
		if len(values) == 0 {
			log.Println("HandleTTSSetVoice: no speakerID selected")
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(discord.Embed{
				Title:       "エラー",
				Description: "話者が選択されていません。もう一度お試しください。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}).WithFlags(discord.MessageFlagEphemeral))
			return err
		}

		speakerID := values[0]
		if !ttsutil.IsSpeakerIDValid(ctx, speakerID) {
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(discord.Embed{
				Title:       "エラー",
				Description: "不正な話者IDが選択されました。もう一度お試しください。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}).WithFlags(discord.MessageFlagEphemeral))
			return err
		}
		memberID := e.Member().User.ID.String()

		memberRepo := repository.NewMemberRepository(ctx.DB)
		if err := memberRepo.Create(memberID); err != nil {
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(discord.Embed{
				Title:       "エラー",
				Description: "メンバー情報の作成に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}).WithFlags(discord.MessageFlagEphemeral))
			return err
		}

		repo := repository.NewTTSPersonalSettingRepository(ctx.DB)
		setting, err := repo.GetByMember(memberID)
		if err != nil {
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(discord.Embed{
				Title:       "エラー",
				Description: "TTS個人設定の取得に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}).WithFlags(discord.MessageFlagEphemeral))
			return err
		}

		if setting == nil {
			defaultSetting := repository.DefaultTTSPersonalSetting
			setting = &defaultSetting
			setting.MemberID = memberID
			setting.SpeakerID = speakerID
			if err := repo.Create(setting); err != nil {
				_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(discord.Embed{
					Title:       "エラー",
					Description: "TTS個人設定の作成に失敗しました。",
					Color:       config.Colors.Error,
					Footer: &discord.EmbedFooter{
						Text:    fmt.Sprintf("Requested by %s", e.User().Username),
						IconURL: e.User().EffectiveAvatarURL(),
					},
					Timestamp: func() *time.Time {
						t := time.Now()
						return &t
					}(),
				}).WithFlags(discord.MessageFlagEphemeral))
				return err
			}
		} else {
			setting.SpeakerID = speakerID
			if err := repo.Update(setting); err != nil {
				_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(discord.Embed{
					Title:       "エラー",
					Description: "TTS個人設定の更新に失敗しました。",
					Color:       config.Colors.Error,
					Footer: &discord.EmbedFooter{
						Text:    fmt.Sprintf("Requested by %s", e.User().Username),
						IconURL: e.User().EffectiveAvatarURL(),
					},
					Timestamp: func() *time.Time {
						t := time.Now()
						return &t
					}(),
				}).WithFlags(discord.MessageFlagEphemeral))
				return err
			}
		}

		label := ttsutil.ResolveSpeakerLabel(ctx, speakerID)
		err = e.UpdateMessage(discord.NewMessageUpdate().WithEmbeds(discord.Embed{
			Title:       "話者設定を更新しました",
			Description: "選択した話者: " + label,
			Color:       config.Colors.Success,
			Footer: &discord.EmbedFooter{
				Text:    fmt.Sprintf("Requested by %s", e.User().Username),
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}))
		return err
	}
}

// HandleTTSSetVoicePage は話者選択のページ送りを処理します
/* TODO: ページ送りを実装
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
						Footer: &discord.EmbedFooter{
							Text:    fmt.Sprintf("Requested by %s", e.User().Username),
							IconURL: e.User().EffectiveAvatarURL(),
						},
						Timestamp: func() *time.Time {
							t := time.Now()
							return &t
						}(),
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
						Footer: &discord.EmbedFooter{
							Text:    fmt.Sprintf("Requested by %s", e.User().Username),
							IconURL: e.User().EffectiveAvatarURL(),
						},
						Timestamp: func() *time.Time {
							t := time.Now()
							return &t
						}(),
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
						Footer: &discord.EmbedFooter{
							Text:    fmt.Sprintf("Requested by %s", e.User().Username),
							IconURL: e.User().EffectiveAvatarURL(),
						},
						Timestamp: func() *time.Time {
							t := time.Now()
							return &t
						}(),
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
*/
