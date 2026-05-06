package tts

import (
	"context"
	"fmt"
	"log"
	"time"

	"unibot/internal"
	"unibot/internal/bot/voice"
	"unibot/internal/model"
	"unibot/internal/repository"

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

		// DB処理
		repo := repository.NewTTSConnectionRepository(ctx.DB)
		ttsConnection, _ := repo.GetByGuildID(guildID.String())

		go func() {
			conCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			// 非同期で実行することで、メインスレッドのイベントループを止めないようにする
			err := conn.Open(conCtx, *userVoiceState.ChannelID, false, true)
			if err != nil {
				log.Printf("Voice connection failed: %v", err)
				return
			}
			log.Println("Voice connection established with DAVE")

			// 読み上げ開始の準備
			channel, ok := e.Client().Caches.Channel(*userVoiceState.ChannelID)

			channelName := "不明なチャンネル"
			if ok {
				channelName = channel.Name()
			}

			player := voice.GetManager().GetOrCreate(guildID.String(), userVoiceState.ChannelID.String(), conn, ctx)

			player.EnqueueText(voice.QueueItem{
				Text:    fmt.Sprintf("%sに、読み上げを接続しました。", channelName),
				Setting: repository.DefaultTTSPersonalSetting,
			})

			if ttsConnection == nil {
				ttsConnection = &model.TTSConnection{
					GuildID:   guildID.String(),
					ChannelID: e.Channel().ID().String(),
				}
				_ = repo.Create(ttsConnection)
			} else {
				ttsConnection.ChannelID = e.Channel().ID().String()
				_ = repo.Update(ttsConnection)
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
