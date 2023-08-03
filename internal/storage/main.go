package storage

type Metadata struct {
	FileHash string `yaml:"file-hash"`
}

var StorageFileExtension = ".dvsfile"
var MetaFileExtension = ".dvsmeta"
