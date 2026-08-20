package logger

import (
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// 出力フォーマット。
const (
	FormatJSON = "json"
	FormatText = "text"
)

// Options は Init に渡すロガーの設定。
//
// Service / Version / Commit / Branch はビルド情報であり、
// このパッケージは unibot/internal を import できない（循環参照になる）ため
// 呼び出し側から渡してもらう。
type Options struct {
	Level     slog.Level
	Format    string
	AddSource bool
	Output    io.Writer

	Service string
	Version string
	Commit  string
	Branch  string
	Env     string

	// Extra は stdout 以外の追加出力先（Discord シンクなど）。
	// Flusher を実装しているものは Shutdown 時に Close される。
	Extra []slog.Handler

	// warnings は設定読み込み時の非致命的な問題。Init が警告ログとして出す。
	warnings []error
}

// ConfigFromEnv は CONFIG_LOG_* 環境変数から設定を読み込む。
// 不正な値はすべて既定値へフォールバックし、警告として記録する。
func ConfigFromEnv() Options {
	o := Options{
		Format:  FormatJSON,
		Output:  os.Stdout,
		Service: "unibot",
		Env:     os.Getenv("CONFIG_LOG_ENV"),
	}

	level, err := ParseLevel(os.Getenv("CONFIG_LOG_LEVEL"))
	if err != nil {
		o.warnings = append(o.warnings, err)
	}
	o.Level = level

	switch strings.ToLower(strings.TrimSpace(os.Getenv("CONFIG_LOG_FORMAT"))) {
	case "", FormatJSON:
		o.Format = FormatJSON
	case FormatText:
		o.Format = FormatText
	default:
		o.warnings = append(o.warnings, errInvalidFormat(os.Getenv("CONFIG_LOG_FORMAT")))
	}

	if v := os.Getenv("CONFIG_LOG_SOURCE"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			o.AddSource = b
		} else {
			o.warnings = append(o.warnings, err)
		}
	}

	return o
}

type invalidFormatError string

func (e invalidFormatError) Error() string {
	return "unknown log format " + strconv.Quote(string(e)) + ", falling back to json"
}

func errInvalidFormat(s string) error { return invalidFormatError(s) }
