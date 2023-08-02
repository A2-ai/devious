package cmd

import (
	"devious/internal/config"
	"devious/internal/git"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zeebo/blake3"
	"golang.org/x/exp/slog"
)

type Metadata struct {
	FileHash string `yaml:"file-hash"`
}

var storageFileExtension = ".dvsfile"
var metaFileExtension = ".dvsmeta"

func runAddCmd(cmd *cobra.Command, args []string) error {
	// Get git dir
	gitDir, err := git.GetRootDir()
	if err != nil {
		return err
	}

	// Get desired file
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	// Load the conf
	conf, err := config.Load()
	if err != nil {
		return err
	}

	// Open source file
	src, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination file
	fileHash := fmt.Sprintf("%x", blake3.Sum256([]byte(filePath)))

	dstPath := filepath.Join(conf.StorageDir, fileHash) + storageFileExtension

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
	copiedBytes, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	slog.Info("Copied file to storage",
		slog.String("filesize", fmt.Sprintf("%1.2f MB", float64(copiedBytes)/1000000)),
		slog.String("from", filePath), slog.String("to", dstPath))

	// Create + write metadata file
	metadataFile, err := os.Create(filePath + metaFileExtension)
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
	ignoreFile, err := os.OpenFile(filepath.Join(gitDir, ".gitignore"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer ignoreFile.Close()

	_, err = ignoreFile.WriteString("\n\n# Devious entry\n/" + filepath.Base(filePath))
	return err
}

func getAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a file to storage",
		RunE:  runAddCmd,
	}

	cmd.Flags().StringP("file", "f", "", "The file to add to storage")
	cmd.MarkFlagFilename("file")
	cmd.MarkFlagRequired("file")

	return cmd
}
