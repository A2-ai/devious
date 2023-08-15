package log

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

var ColorGreen = color.New(color.FgGreen).Sprint
var ColorRed = color.New(color.FgRed).Sprint
var ColorYellow = color.New(color.FgYellow).Sprint
var ColorFaint = color.New(color.Faint).Sprint
var ColorFile = color.New(color.Faint, color.Bold).Sprint

var logOut io.Writer = os.Stdout

func CaptureOutput(f func() error) (string, error) {
	var buf bytes.Buffer

	logOut = &buf
	err := f()
	logOut = os.Stdout
	if err != nil {
		return buf.String(), err
	}

	return buf.String(), nil
}

func Print(args ...any) {
	if JsonLogging {
		return
	}

	logOut.Write([]byte(fmt.Sprintln(args...)))
}

func PrintLogo() {
	Print("🌀 Devious\n")
}

func PrintNotInitialized() {
	Print(ColorRed("✘"), "Devious is not initialized, run", ColorFaint("dvs init <storage-path>"), "to initialize")
	JsonLogger.Issues = append(JsonLogger.Issues, JsonIssue{
		Severity: "error",
		Message:  "devious is not initialized",
	})
}

func OverwritePreviousLine() {
	if JsonLogging {
		return
	}

	logOut.Write([]byte("\033[1A\033[2K"))
}
