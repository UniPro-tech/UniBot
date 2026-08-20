package logger_test

import (
	"context"
	"testing"

	"unibot/internal/logger"

	"github.com/stretchr/testify/assert"
)

func TestWithTraceIsIdempotent(t *testing.T) {
	ctx, first := logger.WithTrace(context.Background())
	ctx, second := logger.WithTrace(ctx)

	assert.Equal(t, first, second)
	assert.Equal(t, first, logger.TraceIDFrom(ctx))
}

func TestWithNewTraceAlwaysReplaces(t *testing.T) {
	ctx, first := logger.WithTrace(context.Background())
	_, second := logger.WithNewTrace(ctx)

	assert.NotEqual(t, first, second)
}

func TestWithRequestAlwaysGeneratesNewID(t *testing.T) {
	ctx, first := logger.WithRequest(context.Background())
	_, second := logger.WithRequest(ctx)

	assert.NotEqual(t, first, second)
}

func TestIDFormat(t *testing.T) {
	trace := logger.NewTraceID()
	request := logger.NewRequestID()

	assert.Len(t, string(trace), 32, "trace_id は W3C trace-id 互換の 32桁")
	assert.Len(t, string(request), 16, "request_id は W3C span-id 互換の 16桁")
	assert.Regexp(t, `^[0-9a-f]{32}$`, string(trace))
	assert.Regexp(t, `^[0-9a-f]{16}$`, string(request))
}

func TestIDsAreAbsentFromEmptyContext(t *testing.T) {
	assert.Empty(t, logger.TraceIDFrom(context.Background()))
	assert.Empty(t, logger.RequestIDFrom(context.Background()))
	assert.Empty(t, logger.TraceIDFrom(nil))
	assert.Empty(t, logger.RequestIDFrom(nil))
}

// 切り離した goroutine へ ID を引き継ぐ契約を固定する。
// context.WithoutCancel は値を保持したまま cancel / deadline のみ切る。
func TestWithoutCancelPreservesIDs(t *testing.T) {
	base, trace := logger.WithTrace(context.Background())
	base, request := logger.WithRequest(base)

	cancellable, cancel := context.WithCancel(base)
	detached := context.WithoutCancel(cancellable)
	cancel()

	assert.NoError(t, detached.Err(), "切り離した context は cancel されない")
	assert.Equal(t, trace, logger.TraceIDFrom(detached))
	assert.Equal(t, request, logger.RequestIDFrom(detached))
}
