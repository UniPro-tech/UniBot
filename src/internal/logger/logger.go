// Package logger は UniBot の構造化ロギング基盤を提供する。
//
// 出力は JSON 1行 = 1レコードで stdout に書き出され、Kubernetes 経由で
// Grafana / Loki に収集される。context 上の trace_id / request_id は
// ContextHandler が自動的に全レコードへ付与する。
//
// 依存ルール: このパッケージは標準ライブラリ以外を import してはならない。
// unibot/internal (config.go) がこのパッケージを使うため、逆向きの import は
// 循環参照になる。
package logger

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
)

// Flusher は終了時にバッファを吐き出す必要のある出力先。
type Flusher interface {
	Close(ctx context.Context) error
}

var (
	levelVar = new(slog.LevelVar)

	current atomic.Pointer[slog.Logger] // Extra を含む本番用ロガー
	plain   atomic.Pointer[slog.Logger] // stdout のみ。再帰防止に使う

	extrasMu sync.Mutex
	extras   []slog.Handler
)

func init() {
	// Init 前に出力されたログも構造化されるよう、既定のロガーを用意しておく。
	levelVar.Set(LevelInfo)
	h := NewContextHandler(newSinkHandler(os.Stdout, FormatJSON, levelVar, false))
	lg := slog.New(h)
	current.Store(lg)
	plain.Store(lg)
}

// Init はロガーを組み立て、slog のデフォルトロガーとして設定する。
func Init(o Options) *slog.Logger {
	if o.Output == nil {
		o.Output = os.Stdout
	}
	if o.Format == "" {
		o.Format = FormatJSON
	}
	levelVar.Set(o.Level)

	stdout := newSinkHandler(o.Output, o.Format, levelVar, o.AddSource)
	static := staticAttrs(o)

	// stdout のみのロガー。Discord シンクを含まないため、
	// シンク自身や disgo のログをここに流しても再帰しない。
	plainLogger := slog.New(NewContextHandler(stdout)).With(static...)
	plain.Store(plainLogger)

	handlers := append([]slog.Handler{stdout}, o.Extra...)
	lg := slog.New(NewContextHandler(NewFanoutHandler(handlers...))).With(static...)
	current.Store(lg)
	slog.SetDefault(lg)

	extrasMu.Lock()
	extras = o.Extra
	extrasMu.Unlock()

	for _, w := range o.warnings {
		lg.Warn("logger configuration fell back to a default", slog.Any("err", w))
	}

	return lg
}

func staticAttrs(o Options) []any {
	attrs := make([]any, 0, 10)
	if o.Service != "" {
		attrs = append(attrs, slog.String("service", o.Service))
	}
	if o.Version != "" {
		attrs = append(attrs, slog.String("version", o.Version))
	}
	if o.Commit != "" {
		attrs = append(attrs, slog.String("commit", o.Commit))
	}
	if o.Branch != "" {
		attrs = append(attrs, slog.String("branch", o.Branch))
	}
	if o.Env != "" {
		attrs = append(attrs, slog.String("env", o.Env))
	}
	return attrs
}

// L は現在のロガーを返す。
func L() *slog.Logger { return current.Load() }

// Plain は stdout のみに出力するロガーを返す。
//
// Discord シンクを含まないため、シンク自身の診断や disgo の内部ログなど
// 「シンクに流すと再帰しうるもの」はすべてこちらを使うこと。
func Plain(component string) *slog.Logger {
	lg := plain.Load()
	if component == "" {
		return lg
	}
	return lg.With(slog.String("component", component))
}

// SetLevel は実行中にログレベルを変更する。
func SetLevel(l slog.Level) { levelVar.Set(l) }

// Level は現在のログレベルを返す。
func Level() slog.Level { return levelVar.Level() }

// Notice は Notice レベルでログを出力する。slog に対応するショートハンドが無いため用意している。
func Notice(ctx context.Context, msg string, args ...any) {
	L().Log(ctx, LevelNotice, msg, args...)
}

// Shutdown は追加出力先のバッファを吐き出して閉じる。
//
// os.Exit は defer を実行しないため、異常終了パスからも明示的に呼ぶこと。
func Shutdown(ctx context.Context) error {
	extrasMu.Lock()
	handlers := extras
	extras = nil
	extrasMu.Unlock()

	var errs []error
	for _, h := range handlers {
		if f, ok := h.(Flusher); ok {
			if err := f.Close(ctx); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}
