package meta

import (
	"os"
	"testing"
)

func TestCreateFile(t *testing.T) {
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

	// Create file
	filePath := tempDir + "/test.txt"
	err = os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create metadata file
	err = Save(Metadata{
		FileHash: "test",
	}, filePath, false)
	if err != nil {
		t.Error(err)
	}

	// Create metadata file again
	err = Save(Metadata{
		FileHash: "test",
	}, filePath, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateLoadFile(t *testing.T) {
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

	// Create file
	filePath := tempDir + "/test.txt"
	err = os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create metadata file
	err = Save(Metadata{
		FileHash: "test",
	}, filePath, false)
	if err != nil {
		t.Error(err)
	}

	// Load metadata file
	metadata, err := Load(filePath)
	if err != nil {
		t.Error(err)
	}

	// Check if the metadata is correct
	if metadata.FileHash != "test" {
		t.Error("Expected file hash to be test")
	}

	// Load metadata file again
	metadata, err = Load(filePath)
	if err != nil {
		t.Error(err)
	}
}
