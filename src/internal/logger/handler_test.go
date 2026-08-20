package logger_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"unibot/internal/logger"

	"github.com/stretchr/testify/assert"
)

func TestContextHandlerAddsIDs(t *testing.T) {
	var buf bytes.Buffer
	lg := logger.Init(logger.Options{Level: logger.LevelInfo, Output: &buf})

	ctx, trace := logger.WithTrace(context.Background())
	ctx, request := logger.WithRequest(ctx)
	lg.InfoContext(ctx, "interaction handled")

	rec := decodeLine(t, buf.Bytes())
	assert.Equal(t, string(trace), rec[logger.KeyTraceID])
	assert.Equal(t, string(request), rec[logger.KeyRequestID])
}

// ID が無い context では、空文字列ではなくキー自体が現れないこと。
func TestContextHandlerOmitsMissingIDs(t *testing.T) {
	var buf bytes.Buffer
	lg := logger.Init(logger.Options{Level: logger.LevelInfo, Output: &buf})

	lg.InfoContext(context.Background(), "bot started")

	rec := decodeLine(t, buf.Bytes())
	assert.NotContains(t, rec, logger.KeyTraceID)
	assert.NotContains(t, rec, logger.KeyRequestID)
}

// WithAttrs で派生させても ID の付与が続くこと。
func TestContextHandlerSurvivesWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	lg := logger.Init(logger.Options{Level: logger.LevelInfo, Output: &buf})

	ctx, trace := logger.WithTrace(context.Background())
	lg.With(slog.String("component", "tts")).InfoContext(ctx, "enqueued")

	rec := decodeLine(t, buf.Bytes())
	assert.Equal(t, "tts", rec["component"])
	assert.Equal(t, string(trace), rec[logger.KeyTraceID])
}

// fanout は子ごとに閾値が違うため、子の Enabled を再確認してから配ること。
func TestFanoutRespectsPerHandlerLevel(t *testing.T) {
	var low, high bytes.Buffer

	lowHandler := slog.NewJSONHandler(&low, &slog.HandlerOptions{Level: logger.LevelDebug})
	highHandler := slog.NewJSONHandler(&high, &slog.HandlerOptions{Level: logger.LevelError})

	lg := slog.New(logger.NewContextHandler(logger.NewFanoutHandler(lowHandler, highHandler)))
	lg.Log(context.Background(), logger.LevelNotice, "bot ready")

	assert.Contains(t, low.String(), "bot ready", "debug 閾値の出力先には届く")
	assert.Empty(t, high.String(), "error 閾値の出力先には届かない")
}

func TestFanoutDeliversToAllEnabledHandlers(t *testing.T) {
	var a, b bytes.Buffer

	lg := slog.New(logger.NewContextHandler(logger.NewFanoutHandler(
		slog.NewJSONHandler(&a, &slog.HandlerOptions{Level: logger.LevelDebug}),
		slog.NewJSONHandler(&b, &slog.HandlerOptions{Level: logger.LevelDebug}),
	)))

	ctx, trace := logger.WithTrace(context.Background())
	lg.ErrorContext(ctx, "boom")

	assert.Contains(t, a.String(), string(trace))
	assert.Contains(t, b.String(), string(trace))
}

// Plain は追加出力先を含まないロガーを返す（再帰防止の要）。
func TestPlainExcludesExtraHandlers(t *testing.T) {
	var stdout, extra bytes.Buffer

	logger.Init(logger.Options{
		Level:  logger.LevelDebug,
		Output: &stdout,
		Extra:  []slog.Handler{slog.NewJSONHandler(&extra, &slog.HandlerOptions{Level: logger.LevelDebug})},
	})

	logger.L().Error("goes everywhere")
	assert.Contains(t, extra.String(), "goes everywhere")

	extra.Reset()
	logger.Plain("logsink").Error("stdout only")

	assert.Contains(t, stdout.String(), "stdout only")
	assert.Empty(t, extra.String(), "Plain の出力は追加出力先へ流れてはならない")
}
