package file

import (
	"os"
	"testing"
)

func TestGetFileHashSanity(t *testing.T) {
	// Create temp dir
	tempFile, err := os.MkdirTemp(".", "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(tempFile)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Create file
	filePath := tempFile + "/lorem.txt"
	err = os.WriteFile(filePath, []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."), 0644)
	if err != nil {
		t.Fatal(err)
	}

	hash, err := GetFileHash(filePath)
	if err != nil {
		t.Error(err)
	}

	if hash != "71fe44583a6268b56139599c293aeb854e5c5a9908eca00105d81ad5e22b7bb6" {
		t.Error("Hash does not match")
	}
}
