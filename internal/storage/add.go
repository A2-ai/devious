package storage

import (
	"dvs/internal/file"
	"dvs/internal/git"
	"dvs/internal/meta"
	"path/filepath"
)

// Adds a file to storage, returning the file hash
func Add(localPath string, storageDir string, gitDir string, dry bool) (hash string, err error) {
	// Create file hash
	fileHash, err := file.GetFileHash(localPath)
	if err != nil {
		return fileHash, err
	}

	dstPath := filepath.Join(storageDir, fileHash) + FileExtension

	// Copy the file to the storage directory
	err = Copy(localPath, dstPath, dry)
	if err != nil {
		return fileHash, err
	}

	// Create + write metadata file
	metadata := meta.Metadata{
		FileHash: fileHash,
	}
	err = meta.CreateFile(metadata, localPath, dry)
	if err != nil {
		return fileHash, err
	}

	// Add file to gitignore
	err = git.AddIgnoreEntry(gitDir, localPath, dry)

	return fileHash, err
}
