package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

// TraceID は Discord のイベント1件（インタラクション・ゲートウェイイベント）に
// 対応する識別子。そのイベントから派生した goroutine やキュー処理にも伝搬させる。
type TraceID string

// RequestID は trace 配下の作業単位1つに対応する識別子。
// 切り離した goroutine や TTS のキュー項目ごとに、同一 trace のまま新規発行する。
type RequestID string

// ログ出力に用いるキー名。Loki のクエリから参照されるため変更しないこと。
const (
	KeyTraceID   = "trace_id"
	KeyRequestID = "request_id"
)

type ctxKey int

const (
	ctxKeyTrace ctxKey = iota
	ctxKeyRequest
)

// NewTraceID は W3C trace-id 互換の 128bit ID（32桁の小文字16進）を生成する。
func NewTraceID() TraceID {
	var b [16]byte
	// crypto/rand.Read はエラーを返さない（返す場合はプロセスが継続不能）。
	_, _ = rand.Read(b[:])
	return TraceID(hex.EncodeToString(b[:]))
}

// NewRequestID は W3C span-id 互換の 64bit ID（16桁の小文字16進）を生成する。
func NewRequestID() RequestID {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return RequestID(hex.EncodeToString(b[:]))
}

// WithTrace は trace_id が未設定の場合のみ新規発行する（冪等）。
// Mux.DefaultContext が発行した ID をミドルウェアが引き継ぐために使う。
func WithTrace(ctx context.Context) (context.Context, TraceID) {
	if id := TraceIDFrom(ctx); id != "" {
		return ctx, id
	}
	return WithNewTrace(ctx)
}

// WithNewTrace は既存の値を無視して常に新しい trace_id を発行する。
func WithNewTrace(ctx context.Context) (context.Context, TraceID) {
	id := NewTraceID()
	return WithTraceID(ctx, id), id
}

// WithTraceID は指定した trace_id を context に載せる。
func WithTraceID(ctx context.Context, id TraceID) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ctxKeyTrace, id)
}

// WithRequest は常に新しい request_id を発行する。trace_id は変更しない。
func WithRequest(ctx context.Context) (context.Context, RequestID) {
	if ctx == nil {
		ctx = context.Background()
	}
	id := NewRequestID()
	return context.WithValue(ctx, ctxKeyRequest, id), id
}

// TraceIDFrom は context から trace_id を取り出す。未設定なら空文字列。
func TraceIDFrom(ctx context.Context) TraceID {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(ctxKeyTrace).(TraceID)
	return id
}

// RequestIDFrom は context から request_id を取り出す。未設定なら空文字列。
func RequestIDFrom(ctx context.Context) RequestID {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(ctxKeyRequest).(RequestID)
	return id
}
