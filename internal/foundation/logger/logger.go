package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
)

var std = newLoggerFromEnv()

const logTimeFormat = "2006-01-02 15:04:05"

func Set(l *slog.Logger) {
	std = l
}

func ReloadFromEnv() {
	std = newLoggerFromEnv()
}

func Info(msg string, args ...any)  { std.Info(msg, args...) }
func Warn(msg string, args ...any)  { std.Warn(msg, args...) }
func Error(msg string, args ...any) { std.Error(msg, args...) }

// Convenience helper when an error is present.
func Err(msg string, err error, args ...any) {
	if err == nil {
		std.Error(msg, args...)
		return
	}
	std.Error(msg, append(args, "error", err)...)
}

//nolint:depguard // The foundation logger may use the standard slog inside the adapter.
func newLoggerFromEnv() *slog.Logger {
	format := strings.ToLower(strings.TrimSpace(env("LOG_FORMAT")))
	level := parseLevel(strings.TrimSpace(env("LOG_LEVEL")))
	colorEnabled := parseLogColor(env("LOG_COLOR"))
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey && attr.Value.Kind() == slog.KindTime {
				return slog.String(attr.Key, attr.Value.Time().Format(logTimeFormat))
			}
			return attr
		},
	}

	if level == levelOff {
		return slog.New(slog.NewJSONHandler(io.Discard, opts))
	}

	switch format {
	case "console":
		return slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level:      level,
			AddSource:  false,
			TimeFormat: logTimeFormat,
			NoColor:    !colorEnabled,
		}))
	default:
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
}

func env(key string) string {
	return os.Getenv(key)
}

var (
	levelDebug slog.Level = slog.LevelDebug
	levelInfo  slog.Level = slog.LevelInfo
	levelError slog.Level = slog.LevelError
	levelOff   slog.Level = slog.Level(127) // effectively disables logging
)

func parseLevel(v string) slog.Leveler {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "debug":
		return levelDebug
	case "error":
		return levelError
	case "off":
		return levelOff
	default:
		return levelInfo
	}
}

func parseLogColor(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "", "1", "true":
		return true
	case "0", "false":
		return false
	default:
		return true
	}
}
