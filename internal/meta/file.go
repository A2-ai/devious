package meta

import (
	"encoding/json"
	"os"
	"time"
)

type Metadata struct {
	FileHash  string    `json:"blake3_checksum"`
	FileSize  uint64    `json:"size"`
	Timestamp time.Time `json:"add_time"`
	Message   string    `json:"message"`
	SavedBy   string    `json:"saved_by"`
}

// Creates a metadata file
func Save(metadata Metadata, path string, dry bool) (err error) {
	var metadataFile *os.File
	if dry {
		return nil
	}

	metadataFile, err = os.Create(path + FileExtension)
	if err != nil {
		return err
	}
	defer metadataFile.Close()

	jsonBytes, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	_, err = metadataFile.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

// Loads the metadata file for a given file, erroring if it doesn't exist or if it can't be decoded
func Load(path string) (metadata Metadata, err error) {
	metadataFile, err := os.Open(path + FileExtension)
	if err != nil {
		return metadata, err
	}

	err = json.NewDecoder(metadataFile).Decode(&metadata)
	if err != nil {
		return metadata, err
	}

	return metadata, nil
}
