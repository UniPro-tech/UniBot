package tts

import (
	"context"
	"time"
	"unibot/internal"
	"unibot/internal/bot/voice"
	"unibot/internal/repository"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func LoadLeaveCommandContext() discord.ApplicationCommandOption {
	return discord.ApplicationCommandOptionSubCommand{
		Name:        "leave",
		Description: "ボイスチャンネルから退出します",
	}
}

func Leave(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config
		guildID := *e.GuildID()

		// 1. BotがVCに参加しているか確認
		conn := e.Client().VoiceManager.GetConn(guildID)
		if conn == nil {
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(),
				discord.NewMessageCreate().WithContent("Botはボイスチャンネルに参加していません。").WithEphemeral(true))
			return err
		}

		// 2. 切断処理
		defer func() {
			closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			conn.Close(closeCtx)
		}()

		// 3. プレイヤーマネージャーやDBの掃除
		mgr := voice.GetManager()
		if player := mgr.Get(guildID.String()); player != nil {
			player.Close()
			mgr.Delete(guildID.String())
		}

		repo := repository.NewTTSConnectionRepository(ctx.DB)
		_ = repo.DeleteByGuildID(guildID.String())

		// 4. 成功レスポンス
		responseEmbed := discord.Embed{
			Title:       "TTSボイスチャンネル退出",
			Description: "ボイスチャンネルから退出しました。",
			Color:       config.Colors.Success,
			Footer: &discord.EmbedFooter{
				Text:    "Requested by " + e.User().Username,
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time { t := time.Now(); return &t }(),
		}

		_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(),
			discord.NewMessageCreate().WithEmbeds(responseEmbed))
		return err
	}
}
