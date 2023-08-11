package storage

import (
	"os"
	"path/filepath"
	"testing"
)

// Gets a file from storage
func TestGetNoLongerInStorage(t *testing.T) {
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

	// Create test file
	err = os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Add file to storage
	hash, err := Add(
		filepath.Join(tempDir, "test.txt"),
		tempDir,
		tempDir,
		"test message",
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Remove file from storage manually
	err = os.Remove(filepath.Join(tempDir, hash) + FileExtension)
	if err != nil {
		t.Fatal(err)
	}

	// Run function
	// should error since file could not be found in storage
	err = Get(
		filepath.Join(tempDir, "test.txt"),
		tempDir,
		tempDir,
		false,
	)
	if err == nil {
		t.Error("Function did not return error")
	}
}
