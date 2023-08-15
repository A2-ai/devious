package storage

import (
	"dvs/internal/file"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/meta"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

// Adds a file to storage, returning the file hash
func Add(localPath string, storageDir string, gitDir string, message string, dry bool) (hash string, err error) {
	// Get file hash
	log.Print("    Getting hash...")

	fileHash, err := file.GetFileHash(localPath)
	if err != nil {
		return fileHash, err
	}

	log.OverwritePreviousLine()
	log.Print("    Getting hash...", log.ColorGreen("✔"))
	log.JsonLogger.Actions = append(log.JsonLogger.Actions, log.JsonAction{
		Action: "got hash",
		Path:   localPath,
	})

	dstPath := filepath.Join(storageDir, fileHash) + FileExtension

	// Copy the file to the storage directory
	// if the destination already exists, skip copying
	_, err = os.Stat(dstPath)
	if os.IsNotExist(err) {
		err = Copy(localPath, dstPath, dry)
		if err != nil {
			return fileHash, err
		}

		log.OverwritePreviousLine()
		log.Print("    Cleaning up...", log.ColorGreen("✔"))
		log.JsonLogger.Actions = append(log.JsonLogger.Actions, log.JsonAction{
			Action: "copied file",
			Path:   localPath,
		})
	} else {
		log.Print("    Copying...", log.ColorGreen("✔ file already up to date"))
	}

	// Get file size
	fileInfo, err := os.Stat(localPath)
	var fileSize uint64
	if err != nil {
		log.Print(log.ColorBold(log.ColorYellow("    !")), "Failed to get file info", log.ColorFaint(err.Error()))
		log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
			Severity: "warning",
			Message:  "failed to get file info",
			Location: localPath,
		})
		fileSize = 0
	} else {
		fileSize = uint64(fileInfo.Size())
	}

	// Get user
	user, err := user.Current()
	var userName string
	if err != nil {
		log.Print(log.ColorBold(log.ColorYellow("    !")), "Failed to get current user", log.ColorFaint(err.Error()))
		log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
			Severity: "warning",
			Message:  "failed to get current user",
			Location: localPath,
		})
	} else {
		userName = user.Username
	}

	// Create + write metadata file
	metadata := meta.Metadata{
		FileHash:  fileHash,
		FileSize:  fileSize,
		Timestamp: time.Now(),
		User:      userName,
		Message:   message,
	}
	err = meta.Save(metadata, localPath, dry)
	if err != nil {
		return fileHash, err
	}

	// Add file to gitignore
	log.Print("    Adding gitignore entry...")

	err = git.AddIgnoreEntry(gitDir, localPath, dry)
	if err != nil {
		return fileHash, err
	}

	log.OverwritePreviousLine()
	log.Print("    Adding gitignore entry...", log.ColorGreen("✔\n"))
	log.JsonLogger.Actions = append(log.JsonLogger.Actions, log.JsonAction{
		Action: "added gitignore entry",
		Path:   localPath,
	})

	return fileHash, nil
}
