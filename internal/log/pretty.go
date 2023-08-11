package log

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var ColorGreen = color.New(color.FgGreen).Sprint
var ColorRed = color.New(color.FgRed).Sprint
var ColorYellow = color.New(color.FgYellow).Sprint
var ColorFaint = color.New(color.Faint).Sprint
var ColorFile = color.New(color.Faint, color.Bold).Sprint

func Print(args ...any) {
	if JsonLogger != nil {
		return
	}

	os.Stdout.Write([]byte(fmt.Sprintln(args...)))
}

func PrintLogo() {
	Print("👺 Devious\n")
}

func ThrowNotInitialized() {
	Print(ColorRed("✘"), "Devious is not initialized, run", ColorFaint("dvs init <storage-path>"), "to initialize")
	os.Exit(0)
}

func OverwritePreviousLine() {
	if JsonLogger != nil {
		return
	}

	os.Stdout.Write([]byte("\033[1A\033[2K"))
}
