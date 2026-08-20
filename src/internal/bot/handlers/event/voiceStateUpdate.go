package event_handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"unibot/internal"
	"unibot/internal/bot/voice"
	"unibot/internal/query"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"gorm.io/gorm"
)

func VoiceStateUpdate(ctx context.Context, bctx *internal.BotContext, e *events.GuildVoiceStateUpdate) {
	client := e.Client()
	vsu := e.VoiceState
	oldVsu := e.OldVoiceState

	// Botの動作は無視
	if e.Member.User.Bot {
		return
	}

	// Bot自身がVCに接続しているか確認
	conn := client.VoiceManager.GetConn(vsu.GuildID)
	if conn == nil {
		return
	}

	botChannelID := getBotChannelID(e)
	if botChannelID == nil {
		return
	}

	// oldVsu や各 ChannelID を安全に取得（キャッシュ未存在時のnilパニック対策）
	var oldChannelID *snowflake.ID
	oldChannelID = oldVsu.ChannelID
	var newChannelID *snowflake.ID
	newChannelID = vsu.ChannelID

	// 同一チャンネル内での状態変更（マイクミュート、画面共有など）は完全に無視
	if newChannelID != nil && oldChannelID != nil && *newChannelID == *oldChannelID {
		return
	}

	// 状態変化のタイプを正確に判定
	changeType := "ignored"
	if newChannelID != nil && *newChannelID == *botChannelID {
		// BotのいるVCに入ってきた（新規入室、または他VCからの移動）
		changeType = "joined"
	} else if oldChannelID != nil && *oldChannelID == *botChannelID {
		// BotのいるVCから出ていった
		if newChannelID == nil {
			changeType = "left"
		} else {
			changeType = "moved" // 他のVCへ移動した
		}
	}

	if changeType == "ignored" {
		return
	}

	// --- 1. 参加処理 ---
	if changeType == "joined" && newChannelID != nil {
		channel, ok := client.Caches.Channel(*newChannelID)
		if !ok {
			return
		}

		text := fmt.Sprintf("%sが %s に参加しました。", e.Member.EffectiveName(), channel.Name())
		vp := voice.GetManager().GetOrCreatePlayer(
			int64(vsu.GuildID),
			int64(*botChannelID),
			conn,
			bctx,
		)
		vp.EnqueueText(voice.QueueItem{
			// 読み上げは非同期に処理されるため、cancel を切って trace のみ引き継ぐ。
			Ctx:  context.WithoutCancel(ctx),
			Text: text,
		})
		return
	}

	// --- 2. 退出・移動処理 ---
	if changeType == "left" || changeType == "moved" {
		// チャンネル内にまだ人間がいるかチェック
		var stillInChannel bool
		states := e.Client().Caches.VoiceStates(conn.GuildID())
		states(func(state discord.VoiceState) bool {
			if state.ChannelID == nil {
				return true
			}
			channelID := int64(*state.ChannelID)
			if state.UserID.String() == client.ID().String() {
				return true
			}
			if state.UserID.String() == e.Member.User.ID.String() {
				return true
			}
			if channelID == int64(*botChannelID) {
				member, ok := getSafeMember(client, conn.GuildID(), state.UserID)
				if ok && !member.User.Bot {
					//stillInChannel = true
					return false
				}
			}
			return true
		})

		// 誰もいなくなった場合（Botの切断処理）
		if !stillInChannel {
			defer func() {
				closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				conn.Close(closeCtx)
			}()

			data, err := query.TtsConnection.Where(query.TtsConnection.GuildID.Eq(int64(vsu.GuildID))).First()
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				slog.ErrorContext(ctx, "failed to load tts connection",
					slog.String("guild_id", vsu.GuildID.String()), slog.Any("err", err))
			}

			mgr := voice.GetManager()
			player := mgr.GetPlayer(int64(vsu.GuildID))
			if player != nil {
				player.Close()
				mgr.DeletePlayer(int64(vsu.GuildID))
			}

			if err == nil && data != nil {
				textChannelID := data.ChannelID
				if _, err := query.TtsConnection.Where(query.TtsConnection.GuildID.Eq(int64(vsu.GuildID))).Delete(); err != nil {
					slog.ErrorContext(ctx, "failed to delete tts connection",
						slog.String("guild_id", vsu.GuildID.String()), slog.Any("err", err))
				}

				embed := discord.Embed{
					Title:       "TTS接続解除",
					Description: "ボイスチャンネルから誰もいなくなったため、TTSの接続を解除しました。",
					Color:       bctx.Config.Colors.Success,
					Timestamp: func() *time.Time {
						t := time.Now()
						return &t
					}(),
				}

				// textChannelID は int64。string(int64) はルーン変換になってしまうため、
				// snowflake.ID へ直接変換する。
				if _, err := client.Rest.CreateMessage(snowflake.ID(textChannelID), discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				}); err != nil {
					slog.ErrorContext(ctx, "failed to send tts disconnect notice",
						slog.String("channel_id", snowflake.ID(textChannelID).String()), slog.Any("err", err))
				}
			}
			return
		}

		// 単なる一人の退出、または移動通知
		if oldChannelID != nil {
			channel, ok := client.Caches.Channel(*oldChannelID)
			if ok {
				var text string
				if changeType == "moved" && newChannelID != nil {
					// 移動先のチャンネル名を取得して「〇〇に移動しました」にする
					toChannel, ok := client.Caches.Channel(*newChannelID)
					if ok {
						text = fmt.Sprintf("%sが %s に移動しました。", e.Member.EffectiveName(), toChannel.Name())
					} else {
						text = fmt.Sprintf("%sが別のチャンネルに移動しました。", e.Member.EffectiveName())
					}
				} else {
					text = fmt.Sprintf("%sが %s から退出しました。", e.Member.EffectiveName(), channel.Name())
				}

				vp := voice.GetManager().GetOrCreatePlayer(int64(vsu.GuildID), int64(*botChannelID), conn, bctx)
				vp.EnqueueText(voice.QueueItem{
					Ctx:  context.WithoutCancel(ctx),
					Text: text,
				})
			}
		}
		return
	}
}

// getBotChannelID と getSafeMember は既存のままでOK
func getBotChannelID(e *events.GuildVoiceStateUpdate) *snowflake.ID {
	vs, ok := e.Client().Caches.VoiceState(e.VoiceState.GuildID, e.Client().ID())
	if !ok || vs.ChannelID == nil {
		return nil
	}
	return vs.ChannelID
}

func getSafeMember(client *bot.Client, guildID, userID snowflake.ID) (*discord.Member, bool) {
	memberRaw, ok := client.Caches.Member(guildID, userID)
	if !ok {
		member, err := client.Rest.GetMember(guildID, userID)
		if err != nil {
			return nil, false
		}
		return member, true
	}
	return &memberRaw, true
}
