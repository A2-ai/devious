package file

import (
	"fmt"
	"os"

	"github.com/zeebo/blake3"
)

func GetFileHash(path string) (hash string, err error) {
	// Hit cache if we can first
	hash, err = GetCachedHash(path)
	if err == nil {
		return hash, nil
	}

	// Read in file bytes
	fileContents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Hash file contents
	hash = fmt.Sprintf("%x", blake3.Sum256(fileContents))
	if err != nil {
		return "", err
	}

	// Cache the hash so we don't have to hash the file again
	err = WriteHashToCache(path, hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}
