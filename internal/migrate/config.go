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

	oldPath := repoDir + "/.dvs.yaml"

	// We're done if .dvs.yaml doesn't exist at the repo root
	_, err = os.Open(oldPath)
	if os.IsNotExist(err) {
		return files, nil
	}

	// Migrate the config file
	if dry {
		return []string{oldPath}, nil
	}
	err = os.Rename(oldPath, repoDir+"/dvs.yaml")
	if err != nil {
		return files, err
	}

	return []string{oldPath}, nil
}
