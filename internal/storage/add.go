package storage

import (
	"dvs/internal/git"
	"dvs/internal/meta"
	"fmt"
	"path/filepath"

	"github.com/zeebo/blake3"
	"golang.org/x/exp/slog"
)

// Adds a file to storage, returning the file hash
func Add(localPath string, storageDir string, gitDir string, dry bool) (hash string, err error) {
	// Create file hash
	fileHash := fmt.Sprintf("%x", blake3.Sum256([]byte(localPath)))
	slog.Debug("Generated file hash", slog.String("hash", fileHash), slog.String("file", localPath))

	dstPath := filepath.Join(storageDir, fileHash) + FileExtension

	// Copy the file to the storage directory
	err = Copy(localPath, dstPath, dry)
	if err != nil {
		return fileHash, err
	}

	// Create + write metadata file
	metadata := meta.Metadata{
		FileHash: fileHash,
	}
	err = meta.CreateFile(metadata, localPath, dry)
	if err != nil {
		return fileHash, err
	}

	// Add file to gitignore
	err = git.AddIgnoreEntry(gitDir, localPath, dry)

	return fileHash, err
}
