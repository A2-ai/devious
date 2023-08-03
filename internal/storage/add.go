package storage

import (
	"devious/internal/config"
	"devious/internal/git"
	"devious/internal/meta"
	"fmt"
	"path/filepath"

	"github.com/zeebo/blake3"
	"golang.org/x/exp/slog"
)

func Add(filePath string, conf config.Config, gitDir string) error {
	// Create file hash
	fileHash := fmt.Sprintf("%x", blake3.Sum256([]byte(filePath)))
	slog.Debug("Created file hash", slog.String("hash", fileHash), slog.String("file-path", filePath))

	dstPath := filepath.Join(conf.StorageDir, fileHash) + StorageFileExtension

	// Copy the file to the storage directory
	Copy(filePath, dstPath, conf)

	// Create + write metadata file
	metadata := meta.Metadata{
		FileHash: fileHash,
	}
	err := meta.CreateFile(metadata, filePath)
	if err != nil {
		return err
	}

	// Add file to gitignore
	err = git.AddToIgnore(gitDir, dstPath)

	return err
}
