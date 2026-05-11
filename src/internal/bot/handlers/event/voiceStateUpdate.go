package event_handlers

import (
	"context"
	"fmt"
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/bot/voice"
	"unibot/internal/repository"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

func VoiceStateUpdate(ctx *internal.BotContext, e *events.GuildVoiceStateUpdate) {
	client := e.Client()
	vsu := e.VoiceState
	oldVsu := e.OldVoiceState

	// Bot自身の接続状況を確認（disgoのVoiceManagerを使用している想定）
	conn := client.VoiceManager.GetConn(vsu.GuildID)
	if conn == nil {
		return
	}

	// Botの動作は無視
	if e.Member.User.Bot {
		return
	}

	// チャンネル内での状態変更（ミュートなど）は無視
	if vsu.ChannelID != nil && oldVsu.ChannelID != nil && *vsu.ChannelID == *oldVsu.ChannelID {
		return
	}

	botChannelID := getBotChannelID(e)
	if botChannelID == "" {
		return
	}

	changeType := "left"
	if oldVsu.ChannelID == nil || *oldVsu.ChannelID != snowflake.MustParse(botChannelID) {
		changeType = "joined"
	} else if vsu.ChannelID != nil && oldVsu.ChannelID != nil {
		changeType = "moved"
	}

	// --- 1. 参加処理 ---
	if changeType == "joined" && vsu.ChannelID != nil && *vsu.ChannelID == snowflake.MustParse(botChannelID) {
		channel, ok := client.Caches.Channel(*vsu.ChannelID)
		if !ok {
			return
		}

		text := fmt.Sprintf("%sが %s に参加しました。", e.Member.EffectiveName(), channel.Name())
		vp := voice.GetManager().GetOrCreate(e.Member.GuildID.String(), botChannelID, conn, ctx)
		vp.EnqueueText(voice.QueueItem{
			Text:    text,
			Setting: repository.DefaultTTSPersonalSetting,
		})
		return
	}

	// --- 2. 退出処理 ---
	if changeType == "left" && *oldVsu.ChannelID == snowflake.MustParse(botChannelID) {
		// チャンネル内にまだ人間がいるかチェック
		var stillInChannel bool
		states := client.Caches.VoiceStates(conn.GuildID())
		states(func(state discord.VoiceState) bool {
			if state.ChannelID == nil {
				return true
			}
			log.Printf("Debug: %s, %s, %s", state.UserID.String(), client.ID().String(), state.ChannelID.String())

			// 1. Bot自身はカウントしない
			if state.UserID.String() == client.ID().String() {
				return true
			}

			// 2. 「今まさに退出/移動したユーザー」本人もカウントしない
			if state.UserID.String() == e.Member.User.ID.String() {
				return true
			}

			// 3. そのユーザーが現在「Botと同じチャンネル」にいるか判定
			if state.ChannelID != nil && state.ChannelID.String() == botChannelID {
				log.Print("Debug: botchanID matched for user ", state.UserID.String())
				// 4. そのユーザーがBotでないことを確認
				member, ok := client.Caches.Member(vsu.GuildID, state.UserID)
				log.Printf("Debug: member cache lookup for user %s, found: %v", state.UserID.String(), ok)
				if ok && !member.User.Bot {
					stillInChannel = true
					return false
				}
			}
			return true
		})

		// 誰もいなくなった場合
		if !stillInChannel {
			defer func() {
				closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				conn.Close(closeCtx)
			}()

			repo := repository.NewTTSConnectionRepository(ctx.DB)
			data, err := repo.GetByGuildID(vsu.GuildID.String())

			mgr := voice.GetManager()
			player := mgr.Get(vsu.GuildID.String())
			if player != nil {
				player.Close()
				mgr.Delete(vsu.GuildID.String())
			}

			if err == nil && data != nil {
				textChannelID := data.ChannelID
				_ = repo.DeleteByGuildID(vsu.GuildID.String())

				embed := discord.Embed{
					Title:       "TTS接続解除",
					Description: "ボイスチャンネルから誰もいなくなったため、TTSの接続を解除しました。",
					Color:       ctx.Config.Colors.Success,
					Timestamp: func() *time.Time {
						t := time.Now()
						return &t
					}(),
				}

				_, _ = client.Rest.CreateMessage(snowflake.MustParse(textChannelID), discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}
			return
		}

		// 単なる一人の退出通知
		channel, ok := client.Caches.Channel(*oldVsu.ChannelID)
		if ok {
			text := fmt.Sprintf("%sが %s から退出しました。", e.Member.EffectiveName(), channel.Name())
			vp := voice.GetManager().GetOrCreate(vsu.GuildID.String(), botChannelID, conn, ctx)
			vp.EnqueueText(voice.QueueItem{
				Text:    text,
				Setting: repository.DefaultTTSPersonalSetting,
			})
		}
		return
	}

	// --- 3. 移動処理 ---
	if changeType == "moved" && oldVsu.ChannelID != nil && *oldVsu.ChannelID == snowflake.MustParse(botChannelID) {
		if vsu.ChannelID == nil {
			return
		}
		channel, ok := client.Caches.Channel(*vsu.ChannelID)
		if ok {
			text := fmt.Sprintf("%sが %s に移動しました。", e.Member.EffectiveName(), channel.Name())
			vp := voice.GetManager().GetOrCreate(vsu.GuildID.String(), botChannelID, conn, ctx)
			vp.EnqueueText(voice.QueueItem{
				Text:    text,
				Setting: repository.DefaultTTSPersonalSetting,
			})
		}
		return
	}
}

// getBotChannelID はBotが現在入っているVCのIDを返します
func getBotChannelID(e *events.GuildVoiceStateUpdate) string {
	vs, ok := e.Client().Caches.VoiceState(e.VoiceState.GuildID, e.Client().ID())
	if !ok || vs.ChannelID == nil {
		return ""
	}
	return vs.ChannelID.String()
}
