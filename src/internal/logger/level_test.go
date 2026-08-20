package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"unibot/internal/logger"

	"github.com/stretchr/testify/assert"
)

func TestLevelName(t *testing.T) {
	tests := []struct {
		level slog.Level
		want  string
	}{
		{logger.LevelDebug, "debug"},
		{logger.LevelInfo, "info"},
		{logger.LevelNotice, "notice"},
		{logger.LevelWarn, "warn"},
		{logger.LevelError, "error"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, logger.LevelName(tt.level))
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		in      string
		want    slog.Level
		wantErr bool
	}{
		{"", logger.LevelInfo, false},
		{"debug", logger.LevelDebug, false},
		{"DEBUG", logger.LevelDebug, false},
		{" info ", logger.LevelInfo, false},
		{"notice", logger.LevelNotice, false},
		{"warn", logger.LevelWarn, false},
		{"warning", logger.LevelWarn, false},
		{"error", logger.LevelError, false},
		{"bogus", logger.LevelInfo, true},
	}

	for _, tt := range tests {
		got, err := logger.ParseLevel(tt.in)
		assert.Equal(t, tt.want, got, "input %q", tt.in)
		if tt.wantErr {
			assert.Error(t, err, "input %q", tt.in)
		} else {
			assert.NoError(t, err, "input %q", tt.in)
		}
	}
}

// Notice レベルが "INFO+2" ではなく "notice" として出力されることを固定する。
func TestNoticeIsRenderedAsNotice(t *testing.T) {
	var buf bytes.Buffer
	lg := logger.Init(logger.Options{
		Level:  logger.LevelDebug,
		Format: logger.FormatJSON,
		Output: &buf,
	})

	lg.Log(context.Background(), logger.LevelNotice, "bot ready")

	rec := decodeLine(t, buf.Bytes())
	assert.Equal(t, "notice", rec["level"])
	assert.NotContains(t, buf.String(), "INFO+2")
}

func TestStaticAttrsArePresent(t *testing.T) {
	var buf bytes.Buffer
	lg := logger.Init(logger.Options{
		Level:   logger.LevelInfo,
		Format:  logger.FormatJSON,
		Output:  &buf,
		Service: "unibot",
		Version: "9.0.0",
		Commit:  "abc1234",
		Branch:  "feat/go",
	})

	lg.Info("hello")

	rec := decodeLine(t, buf.Bytes())
	assert.Equal(t, "unibot", rec["service"])
	assert.Equal(t, "9.0.0", rec["version"])
	assert.Equal(t, "abc1234", rec["commit"])
	assert.Equal(t, "feat/go", rec["branch"])
}

func decodeLine(t *testing.T, b []byte) map[string]any {
	t.Helper()
	line := bytes.TrimSpace(b)
	if i := bytes.IndexByte(line, '\n'); i >= 0 {
		line = line[:i]
	}
	var rec map[string]any
	err := json.Unmarshal(line, &rec)
	assert.NoError(t, err, "output: %s", string(b))
	return rec
}
