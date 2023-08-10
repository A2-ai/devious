package log

import (
	"os"

	"github.com/lmittmann/tint"
	"golang.org/x/exp/slog"
)

func ConfigureGlobalLogger(level slog.Level) {
	opts := &tint.Options{
		Level: level,
		// TimeFormat: time.RFC822,
	}
	handler := tint.NewHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
