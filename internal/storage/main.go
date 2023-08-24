package storage

import "io/fs"

var (
	FileExtension          = ".dvsfile"
	storageDirPermissions  = fs.FileMode(0777)
	storageFilePermissions = fs.FileMode(0666)
)
