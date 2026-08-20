// Package gormslog は gorm のロガーを slog へ橋渡しする。
//
// 既定では SQL 本文は Debug レベルでしか出力されない。
// 本番（CONFIG_LOG_LEVEL=info）では遅いクエリとエラーだけが記録される。
package gormslog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"unibot/internal/logger"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type slogLogger struct {
	cfg Config
}

var (
	_ gormlogger.Interface = (*slogLogger)(nil)
	_ gorm.ParamsFilter    = (*slogLogger)(nil)
)

// New は gorm に渡すロガーを生成する。
func New(cfg Config) gormlogger.Interface {
	cfg.applyDefaults()
	return &slogLogger{cfg: cfg}
}

func (l *slogLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	next := l.cfg
	next.Level = level
	return &slogLogger{cfg: next}
}

// log は出力先を都度解決する。logger.Init より前に New されても
// シンクを取りこぼさないようにするため。
func (l *slogLogger) log() *slog.Logger {
	if l.cfg.Logger != nil {
		return l.cfg.Logger
	}
	return logger.L()
}

// ParamsFilter はプレースホルダのままの SQL を返すために必要。
//
// gorm は db.Logger が gorm.ParamsFilter を実装しているかを型アサートし、
// 実装していない場合はバインド変数を埋め込んだ SQL を組み立てる。
// 辞書の単語・RSS の URL・メッセージ本文がそのままログに乗るため、
// 既定では params を落として "$1, $2" のままにする。
func (l *slogLogger) ParamsFilter(_ context.Context, sql string, params ...any) (string, []any) {
	if l.cfg.ExpandQueryParams {
		return sql, params
	}
	return sql, nil
}

func (l *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.cfg.Level < gormlogger.Info {
		return
	}
	l.log().DebugContext(ctx, fmt.Sprintf(msg, args...), slog.String("component", "gorm"))
}

func (l *slogLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.cfg.Level < gormlogger.Warn {
		return
	}
	l.log().WarnContext(ctx, fmt.Sprintf(msg, args...), slog.String("component", "gorm"))
}

func (l *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.cfg.Level < gormlogger.Error {
		return
	}
	l.log().ErrorContext(ctx, fmt.Sprintf(msg, args...), slog.String("component", "gorm"))
}

func (l *slogLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.cfg.Level <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)

	switch {
	case err != nil && l.cfg.Level >= gormlogger.Error && l.shouldLogError(err):
		sql, rows := fc()
		l.log().ErrorContext(ctx, "sql error", l.attrs(sql, rows, elapsed, slog.Any("err", err))...)

	case l.cfg.SlowThreshold > 0 && elapsed > l.cfg.SlowThreshold && l.cfg.Level >= gormlogger.Warn:
		sql, rows := fc()
		l.log().WarnContext(ctx, "slow sql",
			l.attrs(sql, rows, elapsed, slog.Int64("threshold_ms", l.cfg.SlowThreshold.Milliseconds()))...)

	case l.cfg.Level >= gormlogger.Info:
		// 正常なクエリは Debug。本番で全クエリを吐かないための格下げ。
		sql, rows := fc()
		l.log().DebugContext(ctx, "sql", l.attrs(sql, rows, elapsed)...)
	}
}

func (l *slogLogger) shouldLogError(err error) bool {
	if l.cfg.LogRecordNotFound {
		return true
	}
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (l *slogLogger) attrs(sql string, rows int64, elapsed time.Duration, extra ...slog.Attr) []any {
	attrs := make([]any, 0, len(extra)+4)
	attrs = append(attrs,
		slog.String("component", "gorm"),
		slog.String("sql", sql),
		slog.Int64("duration_ms", elapsed.Milliseconds()),
	)
	if rows >= 0 {
		attrs = append(attrs, slog.Int64("rows", rows))
	}
	for _, a := range extra {
		attrs = append(attrs, a)
	}
	return attrs
}
