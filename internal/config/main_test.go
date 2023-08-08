package config

import (
	"os"
	"testing"
)

func TestSanityCheck(t *testing.T) {
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

	_, err = Read("../../")
	if err != nil {
		t.Error(err)
	}
}
