package storage

import (
	"dvs/internal/log"
	"io"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

type WriteProgress struct {
	bytes int64
	bar   *progressbar.ProgressBar
}

func (wp *WriteProgress) Write(p []byte) (int, error) {
	n := len(p)
	wp.bytes += int64(n)

	wp.bar.Set64(wp.bytes)

	return n, nil
}

// Copies a file from the source path to the destination path
func Copy(srcPath string, destPath string, dry bool) error {
	// Ignore .. and . paths
	if srcPath == ".." || srcPath == "." {
		return os.ErrInvalid
	}

	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get file size
	srcStat, err := srcFile.Stat()
	if err != nil {
		return err
	}
	srcSize := srcStat.Size()

	// Wrap source file in progress reader
	bar := progressbar.DefaultBytes(srcSize, "    Writing file...")
	src := io.TeeReader(srcFile, &WriteProgress{
		bar: bar,
	})

	// Ensure destination exists
	err = os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		return err
	}

	var dst *os.File
	if !dry {
		// Create destination file
		dst, err = os.Create(destPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy the file
		_, err := io.Copy(dst, src)
		if err != nil {
			return err
		}

		log.OverwritePreviousLine()
		log.Print("    Writing file...", log.ColorGreen("✔"))
		log.JsonLogger.Actions = append(log.JsonLogger.Actions, log.JsonAction{
			Action: "wrote file",
			Path:   destPath,
		})
	}

	log.Print("    Cleaning up...")

	// Close the file again so we can catch any errors
	// https://www.joeshaw.org/dont-defer-close-on-writable-files/
	return dst.Close()
}
