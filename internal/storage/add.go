package storage

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/meta"
	"fmt"
	"path/filepath"

	"github.com/zeebo/blake3"
	"golang.org/x/exp/slog"
)

// Adds a file to storage
func Add(path string, conf config.Config, gitDir string, dry bool) error {
	// Create file hash
	fileHash := fmt.Sprintf("%x", blake3.Sum256([]byte(path)))
	slog.Debug("Generated file hash", slog.String("hash", fileHash), slog.String("file", path))

	dstPath := filepath.Join(conf.StorageDir, fileHash) + FileExtension

	// Copy the file to the storage directory
	Copy(path, dstPath, conf, dry)

	// Create + write metadata file
	metadata := meta.Metadata{
		FileHash: fileHash,
	}
	err := meta.CreateFile(metadata, path, dry)
	if err != nil {
		return err
	}

	// Add file to gitignore
	err = git.AddToIgnore(gitDir, dstPath)

	return err
}
