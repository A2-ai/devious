package log

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/lmittmann/tint"
	"golang.org/x/exp/slog"
)

var ColorGreen = color.New(color.FgGreen).Sprint
var ColorRed = color.New(color.FgRed).Sprint
var ColorYellow = color.New(color.FgYellow).Sprint
var ColorFaint = color.New(color.Faint).Sprint
var ColorFile = color.New(color.Faint, color.Bold).Sprint

func ConfigureGlobalLogger(level slog.Level) {
	opts := &tint.Options{
		Level: level,
		// TimeFormat: time.RFC822,
	}
	handler := tint.NewHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func ThrowNotInitialized() {
	RawLog(ColorRed("✘"), "Devious is not initialized, run", ColorFaint("dvs init <storage-path>"), "to initialize")
	os.Exit(0)
}

func RawLog(args ...any) {
	os.Stdout.Write([]byte(fmt.Sprintln(args...)))
}

func PrintLogo() {
	RawLog("👺 Devious\n")
}

func OverwritePreviousLine() {
	os.Stdout.Write([]byte("\033[1A\033[2K"))
}
