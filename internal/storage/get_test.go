package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

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

	// Remove file locally manually
	err = os.Remove(filepath.Join(tempDir, "test.txt"))
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

func TestGetAgainAfterLocalMod(t *testing.T) {
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
	_, err = Add(
		filepath.Join(tempDir, "test.txt"),
		tempDir,
		tempDir,
		"test message",
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Get from storage
	err = Get(
		filepath.Join(tempDir, "test.txt"),
		tempDir,
		tempDir,
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Wait a bit to ensure modification time is different
	time.Sleep(2 * time.Millisecond)

	// Modify the file locally
	err = os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("test2"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Get from storage again
	// File should return to original state
	err = Get(
		filepath.Join(tempDir, "test.txt"),
		tempDir,
		tempDir,
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Check file contents
	data, err := os.ReadFile(filepath.Join(tempDir, "test.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "test" {
		t.Error("File contents did not match, should have been test but got", string(data))
	}
}
