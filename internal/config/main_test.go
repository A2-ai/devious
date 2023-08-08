package config

import (
	"os"
	"testing"
)

func TestReadWriteSanity(t *testing.T) {
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

	err = Write(Config{
		StorageDir: tempDir,
	}, tempDir)
	if err != nil {
		t.Error(err)
	}

	_, err = Read(tempDir)
	if err != nil {
		t.Error(err)
	}
}

func TestReadWriteInvalidDir(t *testing.T) {
	// Ensure write fails
	err := Write(Config{
		StorageDir: "test",
	}, "nonexistent")
	if err == nil {
		t.Error(err)
	}

	// Ensure read fails
	_, err = Read("nonexistent")
	if err == nil {
		t.Error(err)
	}
}
