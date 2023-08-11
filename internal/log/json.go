package log

import (
	"encoding/json"
	"fmt"
	"os"
)

type JsonAction struct {
	Action string `json:"action"`
	Path   string `json:"path"`
}

type JsonFile struct {
	Action string `json:"action"`
	Status string `json:"status"`
}

type JsonIssue struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Location string `json:"path"`
}

type JsonLog struct {
	Actions []JsonAction        `json:"actions"`
	Files   map[string]JsonFile `json:"files"`
	Issues  []JsonIssue         `json:"issues"`
}

var JsonLogger *JsonLog = &JsonLog{
	Files: make(map[string]JsonFile),
}
var JsonLogging bool = false

func Dump() error {
	if !JsonLogging {
		return nil
	}

	jsonBytes, err := json.MarshalIndent(JsonLogger, "", "    ")
	if err != nil {
		return err
	}

	os.Stdout.Write(jsonBytes)
	fmt.Print("\n")

	return nil
}

func DumpAndExit(code int) {
	err := Dump()
	if err != nil {
		fmt.Println("Failed to dump JSON log:", err)
	}

	os.Exit(code)
}
