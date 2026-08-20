package interaction_handler

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"unibot/internal"
	"unibot/internal/logger"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type deferredKey struct{}

// markDeferred はインタラクションに応答を保留した（defer した）ことを記録する。
// エラー時に followup を送るか新規メッセージを送るかの判断に使う。
func markDeferred(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, deferredKey{}, true)
}

func isDeferred(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	deferred, _ := ctx.Value(deferredKey{}).(bool)
	return deferred
}

// ObservabilityMiddleware は全インタラクションに trace_id / request_id を付与し、
// 実行履歴を記録し、エラーと panic を捕捉してユーザーへ応答する。
//
// グローバルミドルウェアの先頭に登録すること。最も外側で動くことにより、
// CreateMasterRecordMiddleware の失敗もここで拾える。
func ObservabilityMiddleware(bctx *internal.BotContext) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) (err error) {
			base := e.Ctx
			if base == nil {
				base = context.Background()
			}
			// trace_id は Mux.DefaultContext が発行済みなら再利用する。
			ctx, traceID := logger.WithTrace(base)
			ctx, _ = logger.WithRequest(ctx)
			e.Ctx = ctx

			kind, path := interactionInfo(e)
			attrs := baseAttrs(e, kind, path)
			start := time.Now()

			defer func() {
				rec := recover()
				if rec == nil {
					return
				}
				slog.ErrorContext(ctx, "interaction panicked", append(attrs,
					slog.Any("panic", rec),
					slog.String("stack", string(debug.Stack())),
					slog.Int64("duration_ms", time.Since(start).Milliseconds()),
					slog.String("outcome", "panic"),
				)...)
				respondError(ctx, bctx, e, traceID)
				// disgo 側で再度ログされないよう、ここで消費する。
				err = nil
			}()

			if err = next(e); err != nil {
				slog.ErrorContext(ctx, "interaction failed", append(attrs,
					slog.Any("err", err),
					slog.Int64("duration_ms", time.Since(start).Milliseconds()),
					slog.String("outcome", "error"),
				)...)
				respondError(ctx, bctx, e, traceID)
				return nil
			}

			// コマンド実行履歴。Info なので Discord へは通知されない。
			slog.InfoContext(ctx, "interaction handled", append(attrs,
				slog.Int64("duration_ms", time.Since(start).Milliseconds()),
				slog.String("outcome", "ok"),
			)...)
			return nil
		}
	}
}

// ErrorSink は Mux のエラーハンドラ。
//
// 通常は ObservabilityMiddleware がエラーを消費するため到達しない。
// グローバルミドルウェアを経由しない経路（NotFound など）のための保険。
func ErrorSink(e *handler.InteractionEvent, err error) {
	ctx := e.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	kind, path := interactionInfo(e)
	slog.ErrorContext(ctx, "unhandled interaction error", append(baseAttrs(e, kind, path),
		slog.Any("err", err),
		slog.String("outcome", "unhandled"),
	)...)
}

// interactionInfo はインタラクションの種別とルーティングパスを返す。
// 判定は disgo の Mux.OnEvent と揃えてある。
func interactionInfo(e *handler.InteractionEvent) (kind string, path string) {
	switch i := e.Interaction.(type) {
	case discord.ApplicationCommandInteraction:
		if data, ok := i.Data.(discord.SlashCommandInteractionData); ok {
			return "slash", data.CommandPath()
		}
		return "command", "/" + i.Data.CommandName()
	case discord.AutocompleteInteraction:
		return "autocomplete", i.Data.CommandPath()
	case discord.ComponentInteraction:
		return "component", i.Data.CustomID()
	case discord.ModalSubmitInteraction:
		return "modal", i.Data.CustomID
	default:
		return "unknown", ""
	}
}

// baseAttrs はコマンド実行履歴に載せる属性を組み立てる。
//
// コマンドのオプション値（data.Options）は絶対に含めないこと。
// /tts dict add の単語や /rss subscribe の URL はユーザー入力そのものであり、
// ログに載せるとセンシティブなデータが流出する。
func baseAttrs(e *handler.InteractionEvent, kind, path string) []any {
	attrs := []any{
		slog.String("event", "interaction"),
		slog.String("kind", kind),
		slog.String("command", path),
		slog.String("interaction_id", e.ID().String()),
		slog.String("user_id", e.User().ID.String()),
		slog.String("channel_id", e.Channel().ID().String()),
	}
	if guildID := e.GuildID(); guildID != nil {
		attrs = append(attrs, slog.String("guild_id", guildID.String()))
	}
	return attrs
}

// respondError はユーザーへ汎用のエラーを返す。
//
// エラー本文は載せず trace_id のみを示す。ユーザーがこの ID を添えて報告すれば、
// 運用者は Grafana 側で該当のログを特定できる。
func respondError(ctx context.Context, bctx *internal.BotContext, e *handler.InteractionEvent, traceID logger.TraceID) {
	message := discord.NewMessageCreate().
		WithEmbeds(errorEmbed(bctx, traceID)).
		WithEphemeral(true)

	var err error
	// 応答を保留済みなら followup、未応答なら新規メッセージで返す。
	if isDeferred(e.Ctx) {
		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), message)
	} else {
		err = e.CreateMessage(message)
	}
	if err != nil {
		slog.WarnContext(ctx, "failed to send error response to user", slog.Any("err", err))
	}
}

func errorEmbed(bctx *internal.BotContext, traceID logger.TraceID) discord.Embed {
	now := time.Now()
	return discord.Embed{
		Title:       "エラーが発生しました",
		Description: "処理中に問題が発生しました。時間をおいて再度お試しください。\n改善しない場合は、下記の ID を添えてサポートへご連絡ください。",
		Color:       bctx.Config.Colors.Error,
		Fields: []discord.EmbedField{
			{Name: "Trace ID", Value: "`" + string(traceID) + "`"},
		},
		Timestamp: &now,
	}
}
