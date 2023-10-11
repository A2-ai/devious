package migrate

import (
	"dvs/internal/config"
	"dvs/internal/git"
	"dvs/internal/storage"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func migrateStorageFile(storageDir string, path string) (modified bool, err error) {
	// Ensure the file has the correct extension
	ext := filepath.Ext(path)
	if ext != "" {
		newPath := strings.TrimSuffix(path, ext)
		err := os.Rename(path, newPath)
		if err != nil {
			return false, err
		}

		path = newPath
		modified = true
	}

	// Ensure the file is not at the root level
	if filepath.Dir(path) == storageDir {
		// Create new directory
		newDir := storage.GetStoragePath(storageDir, filepath.Base(path))
		err := os.MkdirAll(newDir, storage.StorageDirPermissions)
		if err != nil {
			return false, err
		}

		// Move to correct location
		newPath := storage.GetStoragePath(storageDir, filepath.Base(path))
		err = os.Rename(path, newPath)

		modified = true

		if os.IsExist(err) {
			// File already exists, delete the old one
			err = os.Remove(path)
			if err != nil {
				return modified, err
			}
		} else if err != nil {
			return modified, err
		}
	}

	return modified, nil
}

// Migrates storage files to the latest format, returning a list of files that were modified
func MigrateStorageFiles() (modifiedFiles []string, err error) {
	repoDir, _ := git.GetNearestRepoDir(".")
	config, err := config.Read(repoDir)
	if err != nil {
		return nil, err
	}

	// Iterate over all files in the storage directory
	err = filepath.WalkDir(config.StorageDir, func(path string, d fs.DirEntry, _ error) error {
		// Don't migrate the directories
		if d.IsDir() {
			return nil
		}

		// TODO respect gitignore?
		modified, err := migrateStorageFile(config.StorageDir, path)
		if err != nil {
			return err
		}

		if modified {
			modifiedFiles = append(modifiedFiles, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return modifiedFiles, nil
}
