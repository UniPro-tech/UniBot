package discordsink

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"unibot/internal/logger"

	"github.com/disgoorg/snowflake/v2"
)

// 環境変数名。
const (
	EnvChannelID = "CONFIG_LOG_ERROR_CHANNEL_ID"
	EnvLevel     = "CONFIG_LOG_DISCORD_LEVEL"
)

// Config は Discord 通知シンクの設定。
type Config struct {
	// ChannelID が 0 の場合、シンクは完全に無効になる（goroutine も起動しない）。
	ChannelID snowflake.ID

	// MinLevel 未満のレコードは Discord へ送らない。
	// ゼロ値（slog.LevelInfo）は「未指定」とみなし Notice を既定とする。
	// Info をそのまま閾値にすることは想定していない（コマンド実行履歴が
	// Discord に流れてしまうため）。
	MinLevel slog.Level

	// QueueSize は送信待ちレコードの上限。溢れた分は破棄して件数だけ通知する。
	QueueSize int

	// FlushInterval ごとにバッチをまとめて送出する。
	FlushInterval time.Duration

	// MaxEmbeds は1メッセージあたりの Embed 数の上限（Discord の仕様上 10）。
	MaxEmbeds int

	// Self はシンク自身の診断ログの出力先。
	// 再帰を避けるため、必ず logger.Plain を渡すこと。
	Self *slog.Logger
}

// 既定値。
const (
	defaultQueueSize     = 256
	defaultFlushInterval = 5 * time.Second
	defaultMaxEmbeds     = 10
)

// ConfigFromEnv は CONFIG_LOG_* 環境変数から設定を読み込む。
// 不正な値の場合は警告を出して通知を無効化する（プロセスは落とさない）。
func ConfigFromEnv() Config {
	cfg := Config{MinLevel: logger.LevelNotice}
	self := logger.Plain("logsink")

	if raw := strings.TrimSpace(os.Getenv(EnvChannelID)); raw != "" {
		id, err := snowflake.Parse(raw)
		if err != nil {
			self.Warn("invalid log channel id, discord notification disabled",
				slog.String("env", EnvChannelID), slog.Any("err", err))
		} else {
			cfg.ChannelID = id
		}
	}

	if raw := strings.TrimSpace(os.Getenv(EnvLevel)); raw != "" {
		level, err := logger.ParseLevel(raw)
		if err != nil {
			self.Warn("invalid discord log level, falling back to notice",
				slog.String("env", EnvLevel), slog.Any("err", err))
		} else {
			cfg.MinLevel = level
		}
	}

	return cfg
}

func (c *Config) applyDefaults() {
	if c.MinLevel == 0 {
		c.MinLevel = logger.LevelNotice
	}
	if c.QueueSize <= 0 {
		c.QueueSize = defaultQueueSize
	}
	if c.FlushInterval <= 0 {
		c.FlushInterval = defaultFlushInterval
	}
	if c.MaxEmbeds <= 0 || c.MaxEmbeds > defaultMaxEmbeds {
		c.MaxEmbeds = defaultMaxEmbeds
	}
	if c.Self == nil {
		c.Self = logger.Plain("logsink")
	}
}
