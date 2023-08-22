package storage

import (
	"dvs/internal/file"
	"dvs/internal/log"
	"dvs/internal/meta"
	"dvs/internal/utils"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Gets a file from storage
func Get(localPath string, storageDir string, gitDir string, dry bool) error {
	/// Get metadata of specified file
	metadata, err := meta.Load(localPath)
	if os.IsNotExist(err) {
		return errors.New("no metadata")
	} else if err != nil {
		return errors.New("failed to load metadata")
	}

	// Get storage path
	storagePath := filepath.Join(storageDir, metadata.FileHash) + FileExtension

	// Check if file is already present locally
	_, err = os.Stat(localPath)
	var localHash string
	if err == nil {
		// Get local file's hash
		log.Print("    Getting local hash...")

		startTime := time.Now()
		localHash, err = file.GetFileHash(localPath)
		endTime := time.Now()

		log.OverwritePreviousLine()
		if err != nil {
			log.Print("    Getting local hash...", log.ColorBold(log.ColorYellow("!")))
		} else {
			log.Print("    Getting local hash...", log.ColorGreen("✔ in ", utils.FormatDuration(endTime.Sub(startTime))))
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
		log.Print("    Copying...", log.ColorGreen("✔ file already up to date\n"))
	}

	return nil
}
