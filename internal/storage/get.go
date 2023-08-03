package storage

import (
	"devious/internal/config"
	"devious/internal/meta"
	"path/filepath"
)

func Get(localPath string, conf config.Config, gitDir string) error {
	/// Get metadata of specified file
	metadata, err := meta.LoadFile(localPath + meta.FileExtension)
	if err != nil {
		return err
	}

	// Get storage path
	storagePath := filepath.Join(conf.StorageDir, metadata.FileHash) + StorageFileExtension

	// Copy file to destination
	Copy(storagePath, localPath, conf)

	return nil
}
