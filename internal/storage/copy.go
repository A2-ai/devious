package storage

import (
	"devious/internal/config"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	dst, err := os.Create(destPath)

	// Create the directory if it doesn't exist
	// Return if there was an error other than the directory not existing
	if err == os.ErrNotExist {
		err = os.MkdirAll(filepath.Dir(destPath), 0755)
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
		slog.String("from", srcPath), slog.String("to", destPath))

	return nil
}
