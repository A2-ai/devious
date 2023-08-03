package storage

import (
	"devious/internal/config"
	"devious/internal/git"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zeebo/blake3"
	"golang.org/x/exp/slog"
)

type Metadata struct {
	FileHash string `yaml:"file-hash"`
}

var StorageFileExtension = ".dvsfile"
var MetaFileExtension = ".dvsmeta"

func Add(filePath string, conf config.Config, gitDir string) error {
	// Open source file
	src, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination file
	fileHash := fmt.Sprintf("%x", blake3.Sum256([]byte(filePath)))
	slog.Info("Created file hash", slog.String("hash", fileHash))

	dstPath := filepath.Join(conf.StorageDir, fileHash) + StorageFileExtension

	dst, err := os.Create(dstPath)

	// Create the directory if it doesn't exist
	// Return if there was an error other than the directory not existing
	if err == os.ErrNotExist {
		err = os.MkdirAll(filepath.Dir(dstPath), 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	defer dst.Close()

	// Copy the file to the storage directory
	slog.Info("Copying file...")
	copiedBytes, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	slog.Info("Copied file to storage",
		slog.String("filesize", fmt.Sprintf("%1.2f MB", float64(copiedBytes)/1000000)),
		slog.String("from", filePath), slog.String("to", dstPath))

	// Create + write metadata file
	metadataFile, err := os.Create(filePath + MetaFileExtension)
	if err != nil {
		return err
	}
	defer metadataFile.Close()

	err = gob.NewEncoder(metadataFile).Encode(Metadata{
		FileHash: fileHash,
	})
	if err != nil {
		return err
	}

	// Add file to gitignore
	err = git.AddToIgnore(gitDir, dstPath)

	return err
}
