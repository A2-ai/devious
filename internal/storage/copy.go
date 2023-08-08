package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"golang.org/x/exp/slog"
)

type WriteProgress struct {
	bytes uint64
	total string
}

func (wp *WriteProgress) Write(p []byte) (int, error) {
	n := len(p)
	wp.bytes += uint64(n)
	os.Stdout.Write([]byte(fmt.Sprint("\033[K\rWriting file... ", humanize.Bytes(wp.bytes), " out of ", wp.total)))
	return n, nil
}

// Copies a file from the source path to the destination path
func Copy(srcPath string, destPath string, dry bool) error {
	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		slog.Error("Failed to open file", slog.String("path", srcPath))
		return err
	}
	defer srcFile.Close()

	// Get file size
	srcStat, err := srcFile.Stat()
	if err != nil {
		slog.Error("Failed to get file info", slog.String("path", srcPath))
		return err
	}
	srcSize := uint64(srcStat.Size())
	srcSizeHuman := humanize.Bytes(srcSize)

	// Wrap source file in progress reader
	src := io.TeeReader(srcFile, &WriteProgress{
		total: srcSizeHuman,
	})

	// Ensure destination exists
	err = os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		slog.Error("Failed to create directory", slog.String("path", filepath.Dir(destPath)))
		return err
	}

	var dst *os.File
	if !dry {
		// Create destination file
		dst, err = os.Create(destPath)
		if err != nil {
			slog.Error("Failed to create copy destination file", slog.String("path", destPath))
			return err
		}
		defer dst.Close()

		// Copy the file
		_, err := io.Copy(dst, src)
		os.Stdout.Write([]byte("\n"))
		if err != nil {
			slog.Error("Failed to copy file", slog.String("path", srcPath))
			return err
		}
	}

	slog.Info("Copied file",
		slog.String("from", srcPath),
		slog.String("to", destPath),
		slog.String("filesize", srcSizeHuman))

	return nil
}
