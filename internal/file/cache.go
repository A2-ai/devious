package file

import (
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
)

type CacheData struct {
	Hash             string
	CreationTime     time.Time
	ModificationTime time.Time
}

// Returns the hash of the file at the given path, or an error if the file cache is invalid
func GetCachedHash(path string) (string, error) {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// Open the cache
	cachePath := filepath.Join(xdg.CacheHome, "dvs", absPath)
	cacheFile, err := os.Open(cachePath)
	if err != nil {
		return "", err
	}
	defer cacheFile.Close()

	// Read the cache contents
	var cache CacheData
	err = gob.NewDecoder(cacheFile).Decode(&cache)
	if err != nil {
		return "", err
	}

	// Get info about the file
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return "", err
	}

	// Ensure creation time matches, if we can get it
	// if cache.CreationTime != fileInfo.Sys() {
	// 	return "", err
	// }

	// Ensure modification time matches
	if cache.ModificationTime != fileInfo.ModTime() {
		// Remove the cache file
		err = os.Remove(cachePath)
		if err != nil {
			return "", err
		}

		return "", errors.New("file modification time does not match cache (invalidating)")
	}

	// Return the hash
	return cache.Hash, nil
}

// Adds the hash of the file at the given path to the cache
func WriteHashToCache(path string, hash string) error {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	cachePath := filepath.Join(xdg.CacheHome, "dvs", absPath)

	// Create parent directories if they don't exist
	err = os.MkdirAll(filepath.Dir(cachePath), 0755)
	if err != nil {
		return err
	}

	// Open the cache, creating it if it doesn't exist
	cacheFile, err := os.OpenFile(cachePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer cacheFile.Close()

	// Get info about the file
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return err
	}

	// Write the cache contents
	err = gob.NewEncoder(cacheFile).Encode(CacheData{
		Hash: hash,
		// CreationTime:     time.Now(),
		ModificationTime: fileInfo.ModTime(),
	})
	return err
}
