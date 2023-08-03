package storage

import (
	"devious/internal/config"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zeebo/blake3"
	"golang.org/x/exp/slog"
)

func Copy(srcPath string, destPath string, conf config.Config) error {
	// Open source file
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination file
	fileHash := fmt.Sprintf("%x", blake3.Sum256([]byte(destPath)))
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

	// Copy the file
	slog.Info("Copying file...")
	copiedBytes, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	slog.Info("Copied file",
		slog.String("filesize", fmt.Sprintf("%1.2f MB", float64(copiedBytes)/1000000)),
		slog.String("from", srcPath), slog.String("to", dstPath))

	return nil
}
