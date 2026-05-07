package messageComponent

import (
	"fmt"
	"log"
	"strconv"
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
		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdate().WithEmbeds(discord.Embed{
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
func HandleTTSSetVoicePage(ctx *internal.BotContext) func(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
	return func(_ discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		pageIndexString := e.Vars["pageIndex"]
		pageIndex, err := strconv.Atoi(pageIndexString)
		if err != nil {
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(discord.Embed{
				Title:       "エラー",
				Description: "無効なページ番号です。",
				Color:       ctx.Config.Colors.Error,
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

		speakers, err := ttsutil.FetchSpeakers(ctx)
		if err != nil {
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(discord.Embed{
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
			}).WithFlags(discord.MessageFlagEphemeral))
			return err
		}

		pages := ttsutil.BuildSpeakerPages(speakers, ttsutil.SpeakerPageSize)
		if len(pages) == 0 {
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(discord.Embed{
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
			}).WithFlags(discord.MessageFlagEphemeral))
			return err
		}

		memberID := e.Member().User.ID.String()
		currentSpeakerID := ttsutil.GetCurrentSpeakerID(ctx, memberID)
		content, components := ttsutil.BuildVoiceMessage(pageIndex, pages, currentSpeakerID)

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdate().WithEmbeds(discord.Embed{
			Title:       "話者を選択してください",
			Description: content,
			Color:       ctx.Config.Colors.Primary,
			Footer: &discord.EmbedFooter{
				Text:    fmt.Sprintf("Requested by %s", e.User().Username),
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}).WithComponents(components...))
		return err
	}
}
