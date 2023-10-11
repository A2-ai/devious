package storage

import (
	"io/fs"
	"path/filepath"
)

var (
	StorageDirPermissions  = fs.FileMode(0777)
	storageFilePermissions = fs.FileMode(0666)
)

func GetStoragePath(storageDir string, fileHash string) string {
	firstHashSegment := fileHash[:2]
	secondHashSegment := fileHash[2:]
	return filepath.Join(storageDir, firstHashSegment, secondHashSegment)
}
