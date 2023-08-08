package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddIgnoreEntryDuplicateEntry(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp(".", "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Run function
	err = AddIgnoreEntry(
		tempDir,
		"test.txt",
		false,
	)
	if err != nil {
		t.Error("Function returned error", err)
	}

	// Store file contents
	expectedFileContents, err := os.ReadFile(filepath.Join(tempDir, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}

	// Run function again
	err = AddIgnoreEntry(
		tempDir,
		"test.txt",
		false,
	)
	if err != nil {
		t.Error("Function returned error", err)
	}

	// Ensure file contents are unchanged
	newFileContents, err := os.ReadFile(filepath.Join(tempDir, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}

	if string(expectedFileContents) != string(newFileContents) {
		t.Error("Expected file contents to be unchanged")
	}
}
