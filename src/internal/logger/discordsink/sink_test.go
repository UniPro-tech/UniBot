package discordsink_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"sync"
	"testing"
	"time"

	"unibot/internal/logger"
	"unibot/internal/logger/discordsink"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeSender struct {
	mu       sync.Mutex
	messages []discord.MessageCreate
}

func (f *fakeSender) CreateMessage(_ snowflake.ID, m discord.MessageCreate, _ ...rest.RequestOpt) (*discord.Message, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.messages = append(f.messages, m)
	return &discord.Message{}, nil
}

func (f *fakeSender) snapshot() []discord.MessageCreate {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]discord.MessageCreate(nil), f.messages...)
}

// newTestSink はフラッシュをテスト側の Close だけで起こすシンクを作る。
func newTestSink(t *testing.T, cfg discordsink.Config) (*discordsink.Sink, *fakeSender) {
	t.Helper()
	if cfg.FlushInterval == 0 {
		cfg.FlushInterval = time.Hour
	}
	sink := discordsink.New(cfg)
	sender := &fakeSender{}
	sink.Attach(sender)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = sink.Close(ctx)
	})
	return sink, sender
}

func closeSink(t *testing.T, sink *discordsink.Sink) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, sink.Close(ctx))
}

// 「間違ってもエラー本文を載せてはならない」という不変条件を固定するテスト。
func TestSinkNeverLeaksMessageOrAttrs(t *testing.T) {
	sink, sender := newTestSink(t, discordsink.Config{ChannelID: 123})
	lg := slog.New(logger.NewContextHandler(sink))

	ctx, trace := logger.WithTrace(context.Background())
	ctx, request := logger.WithRequest(ctx)

	lg.ErrorContext(ctx, "failed to authenticate with token sk-SUPERSECRET",
		slog.String("email", "user@example.com"),
		slog.String("sql", "SELECT * FROM users WHERE id = '42'"),
		slog.String("content", "ユーザーが入力した秘密のメッセージ"),
	)

	closeSink(t, sink)

	raw, err := json.Marshal(sender.snapshot())
	require.NoError(t, err)
	payload := string(raw)

	assert.NotContains(t, payload, "SUPERSECRET")
	assert.NotContains(t, payload, "user@example.com")
	assert.NotContains(t, payload, "SELECT")
	assert.NotContains(t, payload, "秘密のメッセージ")
	assert.NotContains(t, payload, "failed to authenticate")

	assert.Contains(t, payload, string(trace))
	assert.Contains(t, payload, string(request))
	assert.Contains(t, payload, "error")
}

func TestSinkLevelFloor(t *testing.T) {
	sink, _ := newTestSink(t, discordsink.Config{ChannelID: 123})
	ctx := context.Background()

	assert.False(t, sink.Enabled(ctx, logger.LevelDebug))
	assert.False(t, sink.Enabled(ctx, logger.LevelInfo), "コマンド実行履歴(Info)は Discord に流さない")
	assert.True(t, sink.Enabled(ctx, logger.LevelNotice))
	assert.True(t, sink.Enabled(ctx, logger.LevelWarn))
	assert.True(t, sink.Enabled(ctx, logger.LevelError))
}

func TestSinkDisabledWithoutChannel(t *testing.T) {
	sink := discordsink.New(discordsink.Config{})

	assert.False(t, sink.Enabled(context.Background(), logger.LevelError))
	assert.NoError(t, sink.Handle(context.Background(), slog.Record{Level: logger.LevelError}))
	assert.NoError(t, sink.Close(context.Background()))
	assert.Zero(t, sink.Dropped())
}

func TestSinkDropsWhenQueueFull(t *testing.T) {
	sink, sender := newTestSink(t, discordsink.Config{ChannelID: 123, QueueSize: 2})
	lg := slog.New(logger.NewContextHandler(sink))
	ctx, _ := logger.WithTrace(context.Background())

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 10; i++ {
			lg.ErrorContext(ctx, "boom")
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Handle がブロックした。リクエストパスを止めてはならない")
	}

	assert.EqualValues(t, 8, sink.Dropped())

	closeSink(t, sink)

	// 溢れた件数は overflow Embed で通知される。
	messages := sender.snapshot()
	require.Len(t, messages, 1)
	assert.Contains(t, embedTitles(messages[0]), "🧯 log sink overflow")
}

func TestSinkCoalescesRepeatedRecords(t *testing.T) {
	sink, sender := newTestSink(t, discordsink.Config{ChannelID: 123})
	lg := slog.New(logger.NewContextHandler(sink))

	ctx, _ := logger.WithTrace(context.Background())
	ctx, _ = logger.WithRequest(ctx)
	for i := 0; i < 5; i++ {
		lg.ErrorContext(ctx, "boom")
	}

	closeSink(t, sink)

	messages := sender.snapshot()
	require.Len(t, messages, 1)
	require.Len(t, messages[0].Embeds, 1, "同一 trace/request のレコードは1つに合体する")
	assert.Equal(t, "5", fieldValue(t, messages[0].Embeds[0], "count"))
}

func TestSinkSendsOnlyAllowedFields(t *testing.T) {
	sink, sender := newTestSink(t, discordsink.Config{ChannelID: 123})
	lg := slog.New(logger.NewContextHandler(sink))

	ctx, _ := logger.WithTrace(context.Background())
	ctx, _ = logger.WithRequest(ctx)
	lg.Log(ctx, logger.LevelNotice, "bot ready")

	closeSink(t, sink)

	messages := sender.snapshot()
	require.Len(t, messages, 1)
	require.Len(t, messages[0].Embeds, 1)

	embed := messages[0].Embeds[0]
	assert.Empty(t, embed.Description, "Description は使わない")
	assert.Nil(t, embed.Footer, "Footer は使わない")
	assert.Empty(t, messages[0].Content, "Content は使わない")
	assert.NotNil(t, embed.Timestamp)

	names := make([]string, 0, len(embed.Fields))
	for _, f := range embed.Fields {
		names = append(names, f.Name)
	}
	assert.Equal(t, []string{"level", logger.KeyTraceID, logger.KeyRequestID}, names)
}

func embedTitles(m discord.MessageCreate) []string {
	titles := make([]string, 0, len(m.Embeds))
	for _, e := range m.Embeds {
		titles = append(titles, e.Title)
	}
	return titles
}

func fieldValue(t *testing.T, e discord.Embed, name string) string {
	t.Helper()
	for _, f := range e.Fields {
		if f.Name == name {
			return f.Value
		}
	}
	t.Fatalf("field %q not found", name)
	return ""
}

// Embed の枠に収まらない group を黙って捨てないこと。
func TestSinkCountsGroupsBeyondEmbedLimit(t *testing.T) {
	sink, sender := newTestSink(t, discordsink.Config{ChannelID: 123, MaxEmbeds: 3})
	lg := slog.New(logger.NewContextHandler(sink))

	// 11 個の異なる trace を投入する。
	const total = 11
	for i := 0; i < total; i++ {
		ctx, _ := logger.WithNewTrace(context.Background())
		ctx, _ = logger.WithRequest(ctx)
		lg.ErrorContext(ctx, "boom")
	}

	closeSink(t, sink)

	messages := sender.snapshot()
	require.Len(t, messages, 1)

	embeds := messages[0].Embeds
	assert.LessOrEqual(t, len(embeds), 3, "Discord の Embed 上限を超えない")

	// 最後の Embed は溢れ通知で、載せられなかった件数が出ている。
	last := embeds[len(embeds)-1]
	require.Equal(t, "🧯 log sink overflow", last.Title)
	assert.Equal(t, strconv.Itoa(total-(len(embeds)-1)), fieldValue(t, last, "dropped"))
}

// MaxEmbeds が 1 でも Discord の上限を超えるメッセージを作らないこと。
func TestSinkRespectsMaxEmbedsOfOne(t *testing.T) {
	sink, sender := newTestSink(t, discordsink.Config{ChannelID: 123, MaxEmbeds: 1})
	lg := slog.New(logger.NewContextHandler(sink))

	for i := 0; i < 5; i++ {
		ctx, _ := logger.WithNewTrace(context.Background())
		lg.ErrorContext(ctx, "boom")
	}

	closeSink(t, sink)

	messages := sender.snapshot()
	require.Len(t, messages, 1)
	require.Len(t, messages[0].Embeds, 1)
	assert.Equal(t, "🧯 log sink overflow", messages[0].Embeds[0].Title)
	assert.Equal(t, "5", fieldValue(t, messages[0].Embeds[0], "dropped"))
}

// 再試行で通知が重複しないよう nonce が付いていること。
func TestSinkSetsNonceForIdempotency(t *testing.T) {
	sink, sender := newTestSink(t, discordsink.Config{ChannelID: 123})
	lg := slog.New(logger.NewContextHandler(sink))

	ctx, _ := logger.WithTrace(context.Background())
	lg.ErrorContext(ctx, "boom")

	closeSink(t, sink)

	messages := sender.snapshot()
	require.Len(t, messages, 1)
	assert.NotEmpty(t, messages[0].Nonce)
	assert.True(t, messages[0].EnforceNonce)
}
