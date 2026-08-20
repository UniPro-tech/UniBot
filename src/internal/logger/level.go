package logger

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// syslog に倣ったレベル定義。slog の標準レベルに Notice を足している。
//
// Notice は「運用者が知るべき正常なライフサイクルイベント」を表す。
// Discord への通知はこの Notice 以上を既定の閾値としている。
const (
	LevelDebug  = slog.LevelDebug // -4
	LevelInfo   = slog.LevelInfo  //  0
	LevelNotice = slog.Level(2)   //  Info < Notice < Warn
	LevelWarn   = slog.LevelWarn  //  4
	LevelError  = slog.LevelError //  8
)

// LevelName はレベルを Grafana / Loki が解釈できる小文字表記に変換する。
//
// slog.Level(2).String() は既定で "INFO+2" という表記になってしまうため、
// 独自に名前解決する必要がある。範囲比較ではなく完全一致で解決すること。
func LevelName(l slog.Level) string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelNotice:
		return "notice"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return strings.ToLower(l.String())
	}
}

// ParseLevel は環境変数の文字列をレベルに変換する。
//
// 不正な値であってもプロセスを落とさないよう、必ず利用可能なレベル
// （info）を返した上でエラーを併せて返す。呼び出し側は警告ログに留めること。
func ParseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "notice":
		return LevelNotice, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	default:
		return LevelInfo, fmt.Errorf("unknown log level %q", s)
	}
}

// replaceAttr はトップレベルの level / time 属性のみを整形する。
// groups が空でない場合はネストした属性なので触らない。
func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if len(groups) != 0 {
		return a
	}

	switch a.Key {
	case slog.LevelKey:
		if lvl, ok := a.Value.Any().(slog.Level); ok {
			a.Value = slog.StringValue(LevelName(lvl))
		}
	case slog.TimeKey:
		if a.Value.Kind() == slog.KindTime {
			a.Value = slog.StringValue(a.Value.Time().UTC().Format(time.RFC3339Nano))
		}
	}

	return a
}
