package gormslog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"unibot/internal/logger"
	"unibot/internal/logger/gormslog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func newTestLogger(t *testing.T, cfg gormslog.Config) (gormlogger.Interface, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	cfg.Logger = slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: logger.LevelDebug,
	}))
	return gormslog.New(cfg), &buf
}

func fc(sql string, rows int64) func() (string, int64) {
	return func() (string, int64) { return sql, rows }
}

// 正常なクエリは Info ではなく Debug になる（本番で全クエリを吐かないため）。
func TestSuccessfulQueryIsDebug(t *testing.T) {
	l, buf := newTestLogger(t, gormslog.Config{Level: gormlogger.Info})

	l.Trace(context.Background(), time.Now(), fc("SELECT 1", 1), nil)

	rec := decode(t, buf.Bytes())
	assert.Equal(t, "DEBUG", rec["level"])
	assert.Equal(t, "sql", rec["msg"])
	assert.Equal(t, "SELECT 1", rec["sql"])
}

func TestSuccessfulQuerySilentAtWarnLevel(t *testing.T) {
	l, buf := newTestLogger(t, gormslog.Config{Level: gormlogger.Warn})

	l.Trace(context.Background(), time.Now(), fc("SELECT 1", 1), nil)

	assert.Empty(t, buf.String(), "既定(warn)では正常クエリを出力しない")
}

// ErrRecordNotFound はエラーとして扱わない。
// TTS 未設定のギルドなど正常系で頻出し、Error にすると Discord が溢れるため。
func TestRecordNotFoundIsNotAnError(t *testing.T) {
	l, buf := newTestLogger(t, gormslog.Config{Level: gormlogger.Info})

	l.Trace(context.Background(), time.Now(), fc("SELECT 1", 0), gorm.ErrRecordNotFound)

	rec := decode(t, buf.Bytes())
	assert.Equal(t, "DEBUG", rec["level"])
	assert.Equal(t, "sql", rec["msg"])
	assert.NotContains(t, rec, "err")
}

func TestRecordNotFoundIsSilentAtDefaultLevel(t *testing.T) {
	l, buf := newTestLogger(t, gormslog.Config{})

	l.Trace(context.Background(), time.Now(), fc("SELECT 1", 0), gorm.ErrRecordNotFound)

	assert.Empty(t, buf.String())
}

func TestRecordNotFoundLoggedWhenRequested(t *testing.T) {
	l, buf := newTestLogger(t, gormslog.Config{Level: gormlogger.Error, LogRecordNotFound: true})

	l.Trace(context.Background(), time.Now(), fc("SELECT 1", 0), gorm.ErrRecordNotFound)

	rec := decode(t, buf.Bytes())
	assert.Equal(t, "ERROR", rec["level"])
	assert.Equal(t, "sql error", rec["msg"])
}

func TestSlowQueryIsWarn(t *testing.T) {
	l, buf := newTestLogger(t, gormslog.Config{
		Level:         gormlogger.Warn,
		SlowThreshold: time.Millisecond,
	})

	l.Trace(context.Background(), time.Now().Add(-time.Second), fc("SELECT pg_sleep(1)", 1), nil)

	rec := decode(t, buf.Bytes())
	assert.Equal(t, "WARN", rec["level"])
	assert.Equal(t, "slow sql", rec["msg"])
	assert.EqualValues(t, 1, rec["threshold_ms"])
}

// バインド変数を展開しないことがセキュリティ要件。
func TestParamsFilterDropsParamsByDefault(t *testing.T) {
	l, _ := newTestLogger(t, gormslog.Config{})

	filter, ok := l.(gorm.ParamsFilter)
	require.True(t, ok, "gorm.ParamsFilter を実装していないとユーザー入力が SQL に埋め込まれる")

	sql, params := filter.ParamsFilter(context.Background(),
		"INSERT INTO tts_dictionary (word) VALUES ($1)", "センシティブな単語")

	assert.Equal(t, "INSERT INTO tts_dictionary (word) VALUES ($1)", sql)
	assert.Nil(t, params)
}

func TestParamsFilterExpandsWhenRequested(t *testing.T) {
	l, _ := newTestLogger(t, gormslog.Config{ExpandQueryParams: true})

	filter := l.(gorm.ParamsFilter)
	_, params := filter.ParamsFilter(context.Background(), "SELECT $1", "value")

	assert.Equal(t, []any{"value"}, params)
}

func TestLogModeDoesNotMutateOriginal(t *testing.T) {
	l, buf := newTestLogger(t, gormslog.Config{Level: gormlogger.Warn})

	silent := l.LogMode(gormlogger.Silent)
	silent.Trace(context.Background(), time.Now().Add(-time.Second), fc("SELECT 1", 1), nil)
	assert.Empty(t, buf.String())

	l.Trace(context.Background(), time.Now().Add(-time.Second), fc("SELECT 1", 1), nil)
	assert.Contains(t, buf.String(), "slow sql", "元のロガーの設定は変わらない")
}

func decode(t *testing.T, b []byte) map[string]any {
	t.Helper()
	line := bytes.TrimSpace(b)
	if i := bytes.IndexByte(line, '\n'); i >= 0 {
		line = line[:i]
	}
	var rec map[string]any
	require.NoError(t, json.Unmarshal(line, &rec), "output: %s", string(b))
	return rec
}
