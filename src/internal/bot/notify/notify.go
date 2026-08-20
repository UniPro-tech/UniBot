// Package notify は運用者向けの通知を Discord のチャンネルへ送る。
//
// ログの Discord シンク(internal/logger/discordsink)とは意図的に別経路にしている。
// シンクは「エラー本文を絶対に載せない」ことが不変条件であるのに対し、
// こちらはバージョンやギルド数といった作成済みの内容を載せるため、
// 両者を混同しないようパッケージを分けている。
//
// このパッケージは error 値を受け取らず、整形もしない。
package notify

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"unibot/internal"
	"unibot/internal/logger"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

// 環境変数名。未設定の場合はエラー通知チャンネルにフォールバックする。
const (
	EnvReadyChannelID = "CONFIG_LOG_READY_CHANNEL_ID"
	EnvErrorChannelID = "CONFIG_LOG_ERROR_CHANNEL_ID"
)

// Config は通知先の設定。ChannelID が 0 の場合は通知を行わない。
type Config struct {
	ChannelID snowflake.ID
}

var channelID atomic.Uint64

// ConfigFromEnv は通知先チャンネルを環境変数から読み込む。
func ConfigFromEnv() Config {
	self := logger.Plain("notify")

	for _, env := range []string{EnvReadyChannelID, EnvErrorChannelID} {
		raw := strings.TrimSpace(os.Getenv(env))
		if raw == "" {
			continue
		}
		id, err := snowflake.Parse(raw)
		if err != nil {
			self.Warn("invalid notification channel id",
				slog.String("env", env), slog.Any("err", err))
			continue
		}
		return Config{ChannelID: id}
	}

	return Config{}
}

// Init は通知先を設定する。ChannelID が 0 なら以降の通知はすべて何もしない。
func Init(cfg Config) {
	channelID.Store(uint64(cfg.ChannelID))
}

func target() (snowflake.ID, bool) {
	id := channelID.Load()
	return snowflake.ID(id), id != 0
}

// ReadyInfo は起動通知に載せる情報。
type ReadyInfo struct {
	User   string
	Guilds int
}

// Ready は Bot の起動を通知する。
func Ready(ctx context.Context, client *bot.Client, config *internal.Config, info ReadyInfo) {
	send(ctx, client, discord.Embed{
		Title: "🚀 " + config.BotName + " is ready",
		Color: config.Colors.Success,
		Fields: []discord.EmbedField{
			{Name: "Version", Value: config.BotVersion, Inline: inline()},
			{Name: "Commit", Value: internal.GitCommit, Inline: inline()},
			{Name: "Branch", Value: internal.GitBranch, Inline: inline()},
			{Name: "Guilds", Value: fmt.Sprintf("%d", info.Guilds), Inline: inline()},
			{Name: "User", Value: info.User, Inline: inline()},
		},
		Footer:    traceFooter(ctx),
		Timestamp: now(),
	})
}

// Shutdown は Bot の停止を通知する。
func Shutdown(ctx context.Context, client *bot.Client, config *internal.Config, reason string) {
	send(ctx, client, discord.Embed{
		Title: "🛑 " + config.BotName + " is shutting down",
		Color: config.Colors.Warning,
		Fields: []discord.EmbedField{
			{Name: "Version", Value: config.BotVersion, Inline: inline()},
			{Name: "Reason", Value: reason, Inline: inline()},
		},
		Footer:    traceFooter(ctx),
		Timestamp: now(),
	})
}

func send(ctx context.Context, client *bot.Client, embed discord.Embed) {
	id, ok := target()
	if !ok || client == nil {
		return
	}

	if _, err := client.Rest.CreateMessage(id, discord.MessageCreate{
		Embeds: []discord.Embed{embed},
	}); err != nil {
		// 通知の失敗自体は Warn として記録する。ここはログハンドラではないので再帰しない。
		slog.WarnContext(ctx, "failed to send notification to discord",
			slog.String("channel_id", id.String()), slog.Any("err", err))
	}
}

func traceFooter(ctx context.Context) *discord.EmbedFooter {
	id := logger.TraceIDFrom(ctx)
	if id == "" {
		return nil
	}
	return &discord.EmbedFooter{Text: logger.KeyTraceID + ": " + string(id)}
}

func now() *time.Time {
	t := time.Now()
	return &t
}

func inline() *bool {
	b := true
	return &b
}
