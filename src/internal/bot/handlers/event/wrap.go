package event_handlers

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"unibot/internal"
	"unibot/internal/logger"

	"github.com/disgoorg/disgo/bot"
)

// Wrap は disgo のイベントリスナに trace コンテキスト・panic 復帰・実行時間ログを付与する。
//
// disgo 自体も DispatchEvent で recover しているが、そちらはイベント単位で
// 残りのリスナごと中断され、trace コンテキストも持たない。
// ここではリスナ単位で復帰し、どのイベントで何が起きたかを記録する。
func Wrap[E bot.Event](
	name string,
	bctx *internal.BotContext,
	handle func(ctx context.Context, bctx *internal.BotContext, e E),
) func(e E) {
	return func(e E) {
		ctx, _ := logger.WithNewTrace(context.Background())
		ctx, _ = logger.WithRequest(ctx)

		start := time.Now()

		defer func() {
			if rec := recover(); rec != nil {
				slog.ErrorContext(ctx, "event handler panicked",
					slog.String("event", name),
					slog.Any("panic", rec),
					slog.String("stack", string(debug.Stack())),
					slog.Int64("duration_ms", time.Since(start).Milliseconds()),
					slog.String("outcome", "panic"),
				)
			}
		}()

		handle(ctx, bctx, e)

		// MessageCreate や VoiceStateUpdate は常時発火するため Debug に留める。
		slog.DebugContext(ctx, "event handled",
			slog.String("event", name),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
			slog.String("outcome", "ok"),
		)
	}
}
