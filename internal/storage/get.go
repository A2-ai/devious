package storage

import (
	"dvs/internal/file"
	"dvs/internal/log"
	"dvs/internal/meta"
	"errors"
	"os"
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

	// Get local file's hash
	localHash, _ := file.GetFileHash(localPath)

	// Copy file to destination
	// if the destination already exists and hash matched, skip copying
	_, err = os.Stat(storagePath)
	if os.IsNotExist(err) || metadata.FileHash != localHash {
		err = Copy(storagePath, localPath, dry)
		if err != nil {
			return errors.New("failed to copy file")
		}

		log.OverwritePreviousLine()
		log.Print("    Cleaning up...", log.ColorGreen("✔\n"))
	} else {
		log.Print(log.ColorBold(log.ColorYellow("    !")), "File already exists, not copying\n")
	}

	return nil
}
