package logger

import (
	"context"
	"errors"
	"io"
	"log/slog"
)

// ContextHandler は context 上の trace_id / request_id を全レコードに付与する。
//
// これにより呼び出し側は slog.ErrorContext(ctx, ...) と書くだけでよく、
// 出力先（stdout / Discord）ごとに ID を詰め直す必要がなくなる。
//
// WithGroup された後は付与しない。グループ配下に入れてしまうと
// "db.trace_id" のようにネストしてしまい、Loki 側のクエリが壊れるため。
// 現状このコードベースはグループを使っていない。
type ContextHandler struct {
	inner   slog.Handler
	grouped bool
}

// NewContextHandler は h を包んだ ContextHandler を返す。
// ハンドラチェーンの最も外側に置くこと（1度だけ付与すれば全出力先に届く）。
func NewContextHandler(h slog.Handler) slog.Handler {
	return &ContextHandler{inner: h}
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.grouped {
		if id := TraceIDFrom(ctx); id != "" {
			r.AddAttrs(slog.String(KeyTraceID, string(id)))
		}
		if id := RequestIDFrom(ctx); id != "" {
			r.AddAttrs(slog.String(KeyRequestID, string(id)))
		}
	}
	return h.inner.Handle(ctx, r)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return &ContextHandler{inner: h.inner.WithAttrs(attrs), grouped: h.grouped}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &ContextHandler{inner: h.inner.WithGroup(name), grouped: true}
}

// fanoutHandler は1つのレコードを複数の出力先へ配る。
type fanoutHandler struct {
	handlers []slog.Handler
}

// NewFanoutHandler は複数のハンドラをまとめる。nil は無視する。
func NewFanoutHandler(handlers ...slog.Handler) slog.Handler {
	kept := make([]slog.Handler, 0, len(handlers))
	for _, h := range handlers {
		if h != nil {
			kept = append(kept, h)
		}
	}
	if len(kept) == 1 {
		return kept[0]
	}
	return &fanoutHandler{handlers: kept}
}

func (h *fanoutHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, c := range h.handlers {
		if c.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle は子ごとに閾値が異なる（stdout は info、Discord シンクは notice）ため、
// 必ず子の Enabled を再確認してから Handle を呼ぶ。
// slog 本体は「Enabled が真のときだけ Handle を呼ぶ」契約だが、
// fanout の Enabled は子の OR なので、ここで絞らないと閾値が漏れる。
func (h *fanoutHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for _, c := range h.handlers {
		if !c.Enabled(ctx, r.Level) {
			continue
		}
		// Record は共有できないため子ごとに複製する。
		if err := c.Handle(ctx, r.Clone()); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (h *fanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := make([]slog.Handler, len(h.handlers))
	for i, c := range h.handlers {
		next[i] = c.WithAttrs(attrs)
	}
	return &fanoutHandler{handlers: next}
}

func (h *fanoutHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	next := make([]slog.Handler, len(h.handlers))
	for i, c := range h.handlers {
		next[i] = c.WithGroup(name)
	}
	return &fanoutHandler{handlers: next}
}

// newSinkHandler は stdout へ書き出すハンドラを組み立てる。
// レベルは levelVar 経由で参照するので、実行中の変更が即座に反映される。
func newSinkHandler(w io.Writer, format string, level slog.Leveler, addSource bool) slog.Handler {
	opts := &slog.HandlerOptions{
		Level:       level,
		AddSource:   addSource,
		ReplaceAttr: replaceAttr,
	}
	if format == FormatText {
		return slog.NewTextHandler(w, opts)
	}
	return slog.NewJSONHandler(w, opts)
}
