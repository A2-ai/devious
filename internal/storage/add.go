package storage

import (
	"dvs/internal/file"
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/meta"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

// Adds a file to storage, returning the file hash
func Add(localPath string, storageDir string, gitDir string, dry bool) (hash string, err error) {
	// Create file hash
	if log.JsonLogger == nil {
		log.Print("    Generating hash...")
	}

	fileHash, err := file.GetFileHash(localPath)
	if err != nil {
		return fileHash, err
	}

	if log.JsonLogger == nil {
		log.OverwritePreviousLine()
		log.Print("    Generating hash...", log.ColorGreen("✔"))
	} else {
		log.JsonLogger.Actions = append(log.JsonLogger.Actions, log.JsonAction{
			Action: "generated hash",
			Path:   localPath,
		})
	}

	dstPath := filepath.Join(storageDir, fileHash) + FileExtension

	// Copy the file to the storage directory
	err = Copy(localPath, dstPath, dry)
	if err != nil {
		return fileHash, err
	}

	if log.JsonLogger == nil {
		log.OverwritePreviousLine()
		log.Print("    Cleaning up...", log.ColorGreen("✔"))
	} else {
		log.JsonLogger.Actions = append(log.JsonLogger.Actions, log.JsonAction{
			Action: "copied file",
			Path:   localPath,
		})
	}

	// Get file size
	fileInfo, err := os.Stat(localPath)
	var fileSize uint64
	if err != nil {
		slog.Warn("Failed to get file info", slog.String("path", localPath))
		fileSize = 0
	} else {
		fileSize = uint64(fileInfo.Size())
	}

	// Create + write metadata file
	metadata := meta.Metadata{
		FileHash: fileHash,
		FileSize: fileSize,
	}
	err = meta.Save(metadata, localPath, dry)
	if err != nil {
		return fileHash, err
	}

	// Add file to gitignore
	if log.JsonLogger == nil {
		log.Print("    Adding gitignore entry...")
	} else {
		log.JsonLogger.Actions = append(log.JsonLogger.Actions, log.JsonAction{
			Action: "add gitignore entry",
			Path:   localPath,
		})
	}

	err = git.AddIgnoreEntry(gitDir, localPath, dry)
	if err != nil {
		return fileHash, err
	}

	if log.JsonLogger == nil {
		log.OverwritePreviousLine()
		log.Print("    Adding gitignore entry...", log.ColorGreen("✔"))
	} else {
		log.JsonLogger.Actions = append(log.JsonLogger.Actions, log.JsonAction{
			Action: "add gitignore entry",
			Path:   localPath,
		})
	}
	log.OverwritePreviousLine()
	log.Print("    Adding gitignore entry...", log.ColorGreen("✔\n"))

	return fileHash, nil
}
