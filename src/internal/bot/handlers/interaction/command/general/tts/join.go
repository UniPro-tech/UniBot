package tts

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"unibot/internal"
	"unibot/internal/bot/voice"
	"unibot/internal/model"
	"unibot/internal/query"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func LoadJoinCommandContext() discord.ApplicationCommandOption {
	return discord.ApplicationCommandOptionSubCommand{
		Name:        "join",
		Description: "ボイスチャンネルに参加します",
	}
}

func Join(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config
		guildID := *e.GuildID()

		// ユーザーのボイスステート取得
		userVoiceState, ok := e.Client().Caches.VoiceState(guildID, e.User().ID)

		if !ok || userVoiceState.ChannelID == nil {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "ボイスチャンネルの情報を取得できませんでした。\nボイスチャンネルに参加していますか？",
				Color:       config.Colors.Error,
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

		// Botのボイスステート取得
		botVoiceStatus, botHasVoice := e.Client().Caches.VoiceState(guildID, e.Client().ID())
		if botHasVoice && botVoiceStatus.ChannelID != nil {
			// すでに参加している場合は CreateMessage (Defer済みでない場合) か Update...
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "既にVCに接続しています。",
				Color:       config.Colors.Warning,
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

		// ボイスチャンネル接続
		// disgo の VoiceManager を使用する
		conn := e.Client().VoiceManager.CreateConn(guildID)

		// 接続処理はインタラクションの応答後も続くため、cancel を切って trace のみ引き継ぐ。
		gctx := context.WithoutCancel(e.Ctx)

		go func() {
			conCtx, cancel := context.WithTimeout(gctx, 20*time.Second)
			defer cancel()

			// 非同期で実行することで、メインスレッドのイベントループを止めないようにする
			err := conn.Open(conCtx, *userVoiceState.ChannelID, false, true)
			if err != nil {
				slog.ErrorContext(gctx, "voice connection failed",
					slog.String("guild_id", guildID.String()), slog.Any("err", err))
				notifyJoinFailure(gctx, e, config)
				return
			}
			slog.DebugContext(gctx, "voice connection established with DAVE")

			// 読み上げ開始の準備
			channel, ok := e.Client().Caches.Channel(*userVoiceState.ChannelID)

			channelName := "不明なチャンネル"
			if ok {
				channelName = channel.Name()
			}

			player := voice.GetManager().GetOrCreatePlayer(int64(guildID), int64(*userVoiceState.ChannelID), conn, ctx)

			player.EnqueueText(voice.QueueItem{
				Ctx:  gctx,
				Text: fmt.Sprintf("%sに、読み上げを接続しました。", channelName),
			})

			// DB処理
			ttsConnection := &model.TtsConnection{
				GuildID:   int64(guildID),
				ChannelID: int64(e.Channel().ID()),
			}
			if err := query.TtsConnection.Save(ttsConnection); err != nil {
				slog.ErrorContext(gctx, "failed to save tts connection",
					slog.String("guild_id", guildID.String()), slog.Any("err", err))
			}
		}()

		// 成功レスポンス
		responseEmbed := discord.Embed{
			Title:       "TTSボイスチャンネル接続",
			Description: "ボイスチャンネルに参加しました。",
			Color:       config.Colors.Success,
			Footer: &discord.EmbedFooter{
				Text:    fmt.Sprintf("Requested by %s", e.User().Username),
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}
		_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed).WithEphemeral(false))
		return err
	}
}

// notifyJoinFailure は非同期の接続処理が失敗したことをユーザーへ伝える。
// 成功レスポンスは既に送られているため、追加の followup で訂正する。
func notifyJoinFailure(ctx context.Context, e *handler.CommandEvent, config *internal.Config) {
	embed := discord.Embed{
		Title:       "TTSボイスチャンネル接続",
		Description: "ボイスチャンネルへの接続に失敗しました。もう一度お試しください。",
		Color:       config.Colors.Error,
		Timestamp: func() *time.Time {
			t := time.Now()
			return &t
		}(),
	}
	if _, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(),
		discord.NewMessageCreate().WithEmbeds(embed).WithEphemeral(true)); err != nil {
		slog.WarnContext(ctx, "failed to notify voice connection failure", slog.Any("err", err))
	}
}
