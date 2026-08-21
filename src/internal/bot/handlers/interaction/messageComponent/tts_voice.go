package messageComponent

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"
	"unibot/internal"
	"unibot/internal/bot/ttsutil"
	"unibot/internal/query"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// HandleTTSSetVoice は話者選択のセレクトメニューを処理します
func HandleTTSSetVoice(ctx *internal.BotContext) func(data discord.SelectMenuInteractionData, e *handler.ComponentEvent) error {
	return func(_ discord.SelectMenuInteractionData, e *handler.ComponentEvent) error {
		config := ctx.Config
		values := e.StringSelectMenuInteractionData().Values
		if len(values) == 0 {
			slog.WarnContext(e.Ctx, "no speaker selected", slog.String("outcome", "invalid_input"))
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

		speakerIDstring := values[0]
		speakerIDraw, err := strconv.Atoi(speakerIDstring)
		speakerID := int32(speakerIDraw)
		if !ttsutil.IsSpeakerIDValid(ctx, speakerIDstring) || err != nil {
			if err != nil {
				slog.ErrorContext(e.Ctx, "failed to parse speaker id", slog.Any("err", err))
			}
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
		userID := e.Member().User.ID

		modeValue := e.Vars["global"]
		if modeValue == "true" {
			setting, err := query.TtsUserPreference.Where(query.TtsUserPreference.UserID.Eq(int64(userID))).First()
			if err != nil {
				slog.ErrorContext(e.Ctx, "failed to fetch user tts preference", slog.Any("any", err))
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

			setting.UserID = int64(userID)
			setting.SpeakerID = speakerID
			if err := query.TtsUserPreference.Save(setting); err != nil {
				slog.ErrorContext(e.Ctx, "failed save user tts preference", slog.Any("err", err))
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
		} else {
			{
				setting, err := query.TtsMemberPreference.Where(query.TtsMemberPreference.UserID.Eq(int64(userID))).First()
				if err != nil {
					slog.ErrorContext(e.Ctx, "failed to fetch member tts preference", slog.Any("any", err))
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

				setting.UserID = int64(userID)
				setting.SpeakerID = &speakerID
				if err := query.TtsMemberPreference.Save(setting); err != nil {
					slog.ErrorContext(e.Ctx, "failed save member tts preference", slog.Any("err", err))
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
		}

		label := ttsutil.ResolveSpeakerLabel(ctx, speakerIDstring)
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
		globalModeString := e.Vars["global"]
		isGlobal, err := strconv.ParseBool(globalModeString)
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

		memberID := e.Member().User.ID
		currentSpeakerID, err := ttsutil.GetCurrentSpeakerID(ctx, memberID, isGlobal, e.GuildID())
		if err != nil {
			slog.ErrorContext(e.Ctx, "failed to fetch tts settings", slog.Any("err", err))
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "設定の取得に失敗しました。",
				Color:       ctx.Config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed).WithEphemeral(true))
			return err
		}
		embed, components, err := ttsutil.BuildVoiceMessage(ctx, pageIndex, pages, currentSpeakerID, isGlobal)
		if err != nil {
			slog.ErrorContext(e.Ctx, "failed to generate response message", slog.Any("err", err))
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "メッセージの生成に失敗しました。",
				Color:       ctx.Config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed).WithEphemeral(true))
			return err
		}

		_, err = e.UpdateInteractionResponse(discord.NewMessageUpdate().WithEmbeds(*embed).WithComponents(components...))
		return err
	}
}
