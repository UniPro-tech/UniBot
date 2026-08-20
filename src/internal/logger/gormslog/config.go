package gormslog

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"unibot/internal/logger"

	gormlogger "gorm.io/gorm/logger"
)

// 環境変数名。
const (
	EnvLevel    = "CONFIG_LOG_SQL"
	EnvSlowMsec = "CONFIG_LOG_SQL_SLOW_MS"
)

const defaultSlowThreshold = 200 * time.Millisecond

// Config は gorm ロガーの設定。
//
// bool のフィールドは「true にすると出力が増える」向きで命名している。
// そのため Config{} をそのまま渡しても、SQL にバインド変数を展開せず
// （＝ユーザー入力をログに載せず）、ErrRecordNotFound も記録しない。
//
// 逆向きの命名（ParameterizedQueries / IgnoreRecordNotFound）にすると、
// ゼロ値が「展開する」「記録する」を意味してしまい、フィールドの設定漏れが
// そのままログへの情報漏洩になる。
type Config struct {
	// Level は gorm 側の閾値。ゼロ値は Warn（遅いクエリとエラーのみ）として扱う。
	Level gormlogger.LogLevel

	// SlowThreshold を超えたクエリは Warn として記録する。ゼロ値は 200ms。
	SlowThreshold time.Duration

	// LogRecordNotFound が true のときのみ gorm.ErrRecordNotFound を記録する。
	// TTS 未設定のギルドなど正常系で頻出するため、既定では記録しない。
	LogRecordNotFound bool

	// ExpandQueryParams が true のとき、SQL にバインド変数を展開して出力する。
	// ユーザー入力がそのままログに乗るため、開発時以外は有効にしないこと。
	ExpandQueryParams bool

	// Logger は出力先。nil の場合は logger.L() を都度参照する。
	Logger *slog.Logger
}

// ConfigFromEnv は CONFIG_LOG_SQL 系の環境変数から設定を読み込む。
func ConfigFromEnv() Config {
	cfg := Config{
		Level:         gormlogger.Warn,
		SlowThreshold: defaultSlowThreshold,
	}
	self := logger.Plain("gorm")

	if raw := strings.TrimSpace(os.Getenv(EnvLevel)); raw != "" {
		level, err := parseGormLevel(raw)
		if err != nil {
			self.Warn("invalid sql log level, falling back to warn",
				slog.String("env", EnvLevel), slog.Any("err", err))
		} else {
			cfg.Level = level
		}
	}

	if raw := strings.TrimSpace(os.Getenv(EnvSlowMsec)); raw != "" {
		ms, err := strconv.Atoi(raw)
		if err != nil || ms < 0 {
			self.Warn("invalid slow query threshold, falling back to default",
				slog.String("env", EnvSlowMsec), slog.String("value", raw))
		} else {
			cfg.SlowThreshold = time.Duration(ms) * time.Millisecond
		}
	}

	return cfg
}

func parseGormLevel(s string) (gormlogger.LogLevel, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "silent", "off", "none":
		return gormlogger.Silent, nil
	case "error":
		return gormlogger.Error, nil
	case "warn", "warning":
		return gormlogger.Warn, nil
	case "info", "debug":
		return gormlogger.Info, nil
	default:
		return gormlogger.Warn, unknownLevelError(s)
	}
}

type unknownLevelError string

func (e unknownLevelError) Error() string {
	return "unknown sql log level " + strconv.Quote(string(e))
}

func (c *Config) applyDefaults() {
	if c.Level == 0 {
		c.Level = gormlogger.Warn
	}
	if c.SlowThreshold == 0 {
		c.SlowThreshold = defaultSlowThreshold
	}
}
