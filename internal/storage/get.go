package storage

import (
	"dvs/internal/meta"
	"errors"
	"path/filepath"
)

// Gets a file from storage
func Get(localPath string, storageDir string, gitDir string, dry bool) error {
	/// Get metadata of specified file
	metadata, err := meta.Load(localPath)
	if err != nil {
		return errors.New("failed to load metadata")
	}

	// Get storage path
	storagePath := filepath.Join(storageDir, metadata.FileHash) + FileExtension

	// Copy file to destination
	err = Copy(storagePath, localPath, dry)
	if err != nil {
		return errors.New("failed to copy file")
	}

	return nil
}
