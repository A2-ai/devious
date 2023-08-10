package log

import (
	"encoding/json"
	"os"
)

type JsonAction struct {
	Action string `json:"action"`
	Path   string `json:"path"`
}

type JsonFile struct {
	Path string `json:"path"`
}

type JsonIssue struct {
	Severity string `json:"severity"`
	Path     string `json:"path"`
	Message  string `json:"message"`
}

type JsonLog struct {
	Actions []JsonAction `json:"actions"`
	Files   []JsonFile   `json:"files"`
	Issues  []JsonIssue  `json:"issues"`
}

func Dump(jsonLog *JsonLog) error {
	if jsonLog == nil {
		return nil
	}

	jsonBytes, err := json.MarshalIndent(jsonLog, "", "    ")
	if err != nil {
		return err
	}

	os.Stdout.Write(jsonBytes)

	return nil
}
