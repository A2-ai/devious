package file

import (
	"fmt"
	"os"

	"github.com/zeebo/blake3"
	"golang.org/x/exp/slog"
)

func GetFileHash(path string) (string, error) {
	// Read in file bytes
	fileContents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Hash file contents
	hash := fmt.Sprintf("%x", blake3.Sum256(fileContents))

	slog.Debug("Generated file hash", slog.String("hash", hash), slog.String("file", path))

	return hash, nil
}
