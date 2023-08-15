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

	// Check if file is already present locally
	_, err = os.Stat(localPath)
	var localHash string
	if err == nil {
		// Get local file's hash
		log.Print("    Calculating local hash...")

		localHash, err = file.GetFileHash(localPath)

		log.OverwritePreviousLine()
		if err != nil {
			log.Print("    Calculating local hash...", log.ColorBold(log.ColorYellow("!")))
		} else {
			log.Print("    Calculating local hash...", log.ColorGreen("✔"))
		}
	}

	// Copy file to destination
	// if the destination already exists and hash matched, skip copying
	if os.IsNotExist(err) || metadata.FileHash == "" || localHash == "" || metadata.FileHash != localHash {
		err = Copy(storagePath, localPath, dry)
		if err != nil {
			return errors.New("failed to copy file")
		}

		log.OverwritePreviousLine()
		log.Print("    Cleaning up...", log.ColorGreen("✔\n"))
	} else {
		log.Print(log.ColorGreen("    ✔"), "File already up to date\n")
	}

	return nil
}
