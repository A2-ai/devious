package log

import (
	"os"

	"golang.org/x/exp/slog"
)

func SetGlobalLogLevel(level slog.Level) {
	opts := &slog.HandlerOptions{
		Level: level,
		// see https://www.reddit.com/r/golang/comments/153svuq/slog_how_to_access_the_default_log_format/
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String("tm", a.Value.Time().Format("2006/01/02 15:04:05"))
			}
			if a.Key == slog.LevelKey {
				a.Key = "lvl"
			}
			return a
		},
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
