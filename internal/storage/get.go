package storage

import (
	"dvs/internal/meta"
	"path/filepath"
)

// Gets a file from storage
func Get(localPath string, storageDir string, gitDir string, dry bool) error {
	/// Get metadata of specified file
	metadata, err := meta.LoadFile(localPath)
	if err != nil {
		return err
	}

	// Get storage path
	storagePath := filepath.Join(storageDir, metadata.FileHash) + FileExtension

	// Copy file to destination
	err = Copy(storagePath, localPath, dry)
	if err != nil {
		return err
	}

	return nil
}
