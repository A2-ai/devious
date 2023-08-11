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
	Status string `json:"status"`
}

type JsonIssue struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Location string `json:"path"`
}

type JsonLog struct {
	Actions []JsonAction `json:"actions"`
	Issues  []JsonIssue  `json:"issues"`
}

var JsonLogger *JsonLog = &JsonLog{}
var JsonLogging bool = false

func Dump(s interface{}) error {
	if !JsonLogging {
		return nil
	}

	jsonBytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	os.Stdout.Write(jsonBytes)
	fmt.Print("\n")

	return nil
}

func DumpAndExit(code int) {
	err := Dump(JsonLogger)
	if err != nil {
		fmt.Println("Failed to dump JSON log:", err)
	}

	os.Exit(code)
}
