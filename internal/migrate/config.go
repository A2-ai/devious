package migrate

import (
	"dvs/internal/git"
	"os"
)

func MigrateConfig(dry bool) (files []string, err error) {
	repoDir, err := git.GetNearestRepoDir(".")
	if err != nil {
		return files, err
	}

	// We're done if .dvs.yaml doesn't exist at the repo root
	_, err = os.Open(repoDir)
	if os.IsNotExist(err) {
		return files, nil
	}

	// Migrate te config file
	oldPath := repoDir + "/.dvs.yaml"
	if dry {
		return []string{oldPath}, nil
	}
	os.Rename(oldPath, repoDir+"/dvs.yaml")

	return []string{oldPath}, nil
}
