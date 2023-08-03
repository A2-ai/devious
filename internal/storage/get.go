package storage

import (
	"dvs/internal/config"
	"dvs/internal/meta"
	"path/filepath"
)

// Gets a file from storage
func Get(path string, conf config.Config, gitDir string) error {
	/// Get metadata of specified file
	metadata, err := meta.LoadFile(path)
	if err != nil {
		return err
	}

	// Get storage path
	storagePath := filepath.Join(conf.StorageDir, metadata.FileHash) + FileExtension

	// Copy file to destination
	Copy(storagePath, path, conf)

	return nil
}
