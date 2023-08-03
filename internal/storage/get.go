package storage

import (
	"devious/internal/config"
	"encoding/gob"
	"os"
	"path/filepath"
)

func Get(localPath string, conf config.Config, gitDir string) error {
	/// Get metadata of specified file
	metadataFile, err := os.Open(localPath + MetaFileExtension)
	if err != nil {
		return err
	}

	// Decode metadata
	var metadata Metadata
	err = gob.NewDecoder(metadataFile).Decode(&metadata)
	if err != nil {
		return err
	}

	// Get storage path
	storagePath := filepath.Join(conf.StorageDir, metadata.FileHash) + StorageFileExtension

	// Copy file to destination
	Copy(storagePath, localPath, conf)

	return nil
}
