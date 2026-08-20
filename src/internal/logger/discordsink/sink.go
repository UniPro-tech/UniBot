// Package discordsink は Error / Warn / Notice のログを Discord のチャンネルへ通知する
// slog.Handler を提供する。
//
// # 秘匿に関する不変条件
//
// このパッケージは trace_id / request_id / level / timestamp 以外を
// Discord に送信してはならない。
//
// ログ本文（slog.Record.Message）と属性（Record.Attrs）にはユーザーの
// メッセージ内容・辞書の単語・SQL・URL などが含まれうる。サニタイズが
// 不十分であった場合にセンシティブなデータが露出するため、この実装は
// 本文も属性も一切読まない。ID は Record ではなく context から取得する。
//
// 詳細は stdout（Grafana / Loki）側で trace_id を引いて確認する運用とする。
package discordsink

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"unibot/internal/logger"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
)

// MessageSender は Discord へメッセージを送るための最小限のインターフェース。
// disgo の *bot.Client の Rest フィールド（rest.Rest）がこれを満たす。
type MessageSender interface {
	CreateMessage(channelID snowflake.ID, messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error)
}

// Embed の色。logger 配下は unibot/internal を import できないため定数で持つ。
// 値は internal.Colors と揃えてある。
const (
	colorNotice = 0x3498DB
	colorWarn   = 0xF1C40F
	colorError  = 0xE74C3C
)

const (
	// 1回のフラッシュで扱うレコード数の上限。
	maxBatch = 100
	// 通常時の送信リトライ回数。
	sendAttempts = 3
	// 終了時のフラッシュに使うタイムアウト。
	shutdownFlushTimeout = 5 * time.Second
)

// entry は Discord へ送出しうる唯一の情報。
// ここに本文や属性を足してはならない（パッケージコメント参照）。
type entry struct {
	Trace   logger.TraceID
	Request logger.RequestID
	Level   slog.Level
	At      time.Time
}

// Sink は slog.Handler かつ logger.Flusher。
type Sink struct {
	cfg      Config
	disabled bool

	ch   chan entry
	done chan struct{}
	wg   sync.WaitGroup
	once sync.Once

	// closedMu は closed の確認とキュー投入を同期する。
	// これが無いと「Handle が closed を確認 → Close が worker を終了 →
	// Handle が s.ch に投入」の順で、受理したレコードが失われる。
	closedMu sync.RWMutex
	closed   bool

	senderMu sync.RWMutex
	sender   MessageSender

	dropped atomic.Uint64
}

var (
	_ slog.Handler   = (*Sink)(nil)
	_ logger.Flusher = (*Sink)(nil)
)

// New はシンクを生成する。ChannelID が 0 の場合は何もしない無効なシンクを返す
// （goroutine を起動しないため、呼び出し側は分岐せず常に生成してよい）。
func New(cfg Config) *Sink {
	cfg.applyDefaults()

	if cfg.ChannelID == 0 {
		return &Sink{cfg: cfg, disabled: true}
	}

	s := &Sink{
		cfg:  cfg,
		ch:   make(chan entry, cfg.QueueSize),
		done: make(chan struct{}),
	}
	s.wg.Add(1)
	go s.run()
	return s
}

// Attach は disgo クライアント構築後に送信経路を渡す。
// Attach 前のレコードはキューに保持され、Attach 後の最初のフラッシュで
// 本来のタイムスタンプのまま送出される。
func (s *Sink) Attach(sender MessageSender) {
	if s.disabled {
		return
	}
	s.senderMu.Lock()
	s.sender = sender
	s.senderMu.Unlock()
}

// Dropped はキュー溢れなどで破棄したレコード数の累計を返す（診断用）。
func (s *Sink) Dropped() uint64 { return s.dropped.Load() }

func (s *Sink) Enabled(_ context.Context, level slog.Level) bool {
	if s.disabled || level < s.cfg.MinLevel {
		return false
	}
	s.closedMu.RLock()
	defer s.closedMu.RUnlock()
	return !s.closed
}

// Handle はレコードから ID・レベル・時刻だけを取り出してキューへ積む。
//
// リクエストパスを絶対に止めないため、キューが満杯なら破棄する（ノンブロッキング）。
// 戻り値は常に nil。非 nil を返すと slog が stderr に生ログを吐いてしまう。
func (s *Sink) Handle(ctx context.Context, r slog.Record) error {
	if s.disabled {
		return nil
	}

	// r.Message と r.Attrs は意図的に読まない。ID は context からのみ取得する。
	e := entry{
		Trace:   logger.TraceIDFrom(ctx),
		Request: logger.RequestIDFrom(ctx),
		Level:   r.Level,
		At:      r.Time,
	}
	if e.At.IsZero() {
		e.At = time.Now()
	}

	// closed の確認と投入を不可分にする。RLock なので複数の Handle は並行できる。
	s.closedMu.RLock()
	defer s.closedMu.RUnlock()
	if s.closed {
		return nil
	}

	select {
	case s.ch <- e:
	default:
		s.dropped.Add(1)
	}
	return nil
}

// WithAttrs は属性を意図的に破棄する。Discord へ送るのは ID とレベルのみ。
func (s *Sink) WithAttrs(_ []slog.Attr) slog.Handler { return s }

// WithGroup も同様に何もしない。
func (s *Sink) WithGroup(_ string) slog.Handler { return s }

// Close はキューを吐き出して worker を停止する。
func (s *Sink) Close(ctx context.Context) error {
	if s.disabled {
		return nil
	}
	s.once.Do(func() {
		// 進行中の Handle が投入を終えるまで待ってから worker を止める。
		s.closedMu.Lock()
		s.closed = true
		s.closedMu.Unlock()
		close(s.done)
	})

	waited := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(waited)
	}()

	select {
	case <-waited:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Sink) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.cfg.FlushInterval)
	defer ticker.Stop()

	// pending は「送信できずに持ち越したレコード」のみを保持する。
	// 通常はチャネル自体がバッファであり、溢れた分は Handle 側で破棄される。
	var pending []entry
	for {
		select {
		case <-s.done:
			pending = drain(pending, s.ch)
			ctx, cancel := context.WithTimeout(context.Background(), shutdownFlushTimeout)
			s.flush(ctx, pending, 1)
			cancel()
			return

		case <-ticker.C:
			pending = drain(pending, s.ch)
			pending = s.tryFlush(pending)
		}
	}
}

// tryFlush は送信可能ならバッチを送出して空にする。
// Attach 前は送れないのでバッチに保持したままにする（上限を超えた古い分は破棄）。
func (s *Sink) tryFlush(batch []entry) []entry {
	if len(batch) == 0 && s.dropped.Load() == 0 {
		return batch
	}
	if s.currentSender() == nil {
		if len(batch) > maxBatch {
			s.dropped.Add(uint64(len(batch) - maxBatch))
			batch = append(batch[:0], batch[len(batch)-maxBatch:]...)
		}
		return batch
	}
	s.flush(context.Background(), batch, sendAttempts)
	return nil
}

func (s *Sink) flush(ctx context.Context, batch []entry, attempts int) {
	dropped := s.dropped.Swap(0)
	if len(batch) == 0 && dropped == 0 {
		return
	}

	sender := s.currentSender()
	if sender == nil {
		// 送れなかった分はドロップ扱いに戻す。
		s.dropped.Add(dropped + uint64(len(batch)))
		return
	}

	msg := buildMessage(batch, dropped, s.cfg.MaxEmbeds)
	if len(msg.Embeds) == 0 {
		return
	}

	// CreateMessage はメッセージ作成後にエラーを返すことがある。
	// 再試行で通知が重複しないよう、バッチごとに固定の nonce を付ける。
	msg = msg.WithNonce(string(logger.NewRequestID())).WithEnforceNonce(true)

	var lastErr error
	backoff := time.Second
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				lastErr = ctx.Err()
				attempt = attempts
				continue
			}
			backoff *= 2
		}
		if _, err := sender.CreateMessage(s.cfg.ChannelID, msg); err != nil {
			lastErr = err
			continue
		}
		return
	}

	// 失敗の記録は stdout のみへ。ここでシンク付きロガーを使うと再帰する。
	s.cfg.Self.Warn("failed to deliver log notification to discord",
		slog.Any("err", lastErr),
		slog.Int("embeds", len(msg.Embeds)))
	s.dropped.Add(uint64(len(batch)))
}

func (s *Sink) currentSender() MessageSender {
	s.senderMu.RLock()
	defer s.senderMu.RUnlock()
	return s.sender
}

// drain はチャネルに残っているレコードをすべて回収する。
func drain(batch []entry, ch <-chan entry) []entry {
	for {
		select {
		case e := <-ch:
			batch = append(batch, e)
		default:
			return batch
		}
	}
}

type group struct {
	entry entry
	count int
}

// coalesce は (level, trace_id, request_id) が同じレコードをまとめる。
// エラーストーム時に1つのインタラクションが Embed を埋め尽くすのを防ぐ。
func coalesce(batch []entry) []group {
	type key struct {
		level   slog.Level
		trace   logger.TraceID
		request logger.RequestID
	}

	index := make(map[key]int, len(batch))
	groups := make([]group, 0, len(batch))

	for _, e := range batch {
		k := key{level: e.Level, trace: e.Trace, request: e.Request}
		if i, ok := index[k]; ok {
			groups[i].count++
			continue
		}
		index[k] = len(groups)
		groups = append(groups, group{entry: e, count: 1})
	}

	return groups
}

func buildMessage(batch []entry, dropped uint64, maxEmbeds int) discord.MessageCreate {
	if maxEmbeds <= 0 || maxEmbeds > defaultMaxEmbeds {
		maxEmbeds = defaultMaxEmbeds
	}

	groups := coalesce(batch)

	limit := maxEmbeds
	if dropped > 0 || len(groups) > limit {
		// 溢れ通知の分を1枠空けておく。
		limit = maxEmbeds - 1
	}
	if limit < 0 {
		limit = 0
	}

	// 枠に収まらない group は黙って捨てず、破棄件数として計上する。
	if len(groups) > limit {
		for _, g := range groups[limit:] {
			dropped += uint64(g.count)
		}
		groups = groups[:limit]
	}

	embeds := make([]discord.Embed, 0, len(groups)+1)
	for _, g := range groups {
		embeds = append(embeds, buildEmbed(g))
	}
	if dropped > 0 {
		embeds = append(embeds, overflowEmbed(dropped))
	}

	return discord.MessageCreate{Embeds: embeds}
}

func buildEmbed(g group) discord.Embed {
	level := logger.LevelName(g.entry.Level)

	fields := []discord.EmbedField{
		{Name: "level", Value: level, Inline: boolPtr(true)},
		{Name: logger.KeyTraceID, Value: code(string(g.entry.Trace)), Inline: boolPtr(true)},
		{Name: logger.KeyRequestID, Value: code(string(g.entry.Request)), Inline: boolPtr(true)},
	}
	if g.count > 1 {
		fields = append(fields, discord.EmbedField{
			Name: "count", Value: strconv.Itoa(g.count), Inline: boolPtr(true),
		})
	}

	at := g.entry.At
	return discord.Embed{
		Title:     titleFor(g.entry.Level),
		Color:     colorFor(g.entry.Level),
		Fields:    fields,
		Timestamp: &at,
	}
}

func overflowEmbed(dropped uint64) discord.Embed {
	now := time.Now()
	return discord.Embed{
		Title: "🧯 log sink overflow",
		Color: colorWarn,
		Fields: []discord.EmbedField{
			{Name: "dropped", Value: strconv.FormatUint(dropped, 10), Inline: boolPtr(true)},
		},
		Timestamp: &now,
	}
}

func titleFor(level slog.Level) string {
	switch {
	case level >= logger.LevelError:
		return "🔴 ERROR"
	case level >= logger.LevelWarn:
		return "⚠️ WARN"
	default:
		return "📣 NOTICE"
	}
}

func colorFor(level slog.Level) int {
	switch {
	case level >= logger.LevelError:
		return colorError
	case level >= logger.LevelWarn:
		return colorWarn
	default:
		return colorNotice
	}
}

func code(s string) string {
	if s == "" {
		return "-"
	}
	return "`" + s + "`"
}

func boolPtr(b bool) *bool { return &b }
