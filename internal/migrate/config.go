package migrate

import (
	"dvs/internal/git"
	"os"
)

func MigrateConfig(dry bool) (match bool, err error) {
	repoDir, err := git.GetNearestRepoDir(".")
	if err != nil {
		return false, err
	}

	// We're done if .dvs.yaml doesn't exist at the repo root
	_, err = os.Open(repoDir)
	if os.IsNotExist(err) {
		return false, nil
	}

	// Migrate te config file
	if dry {
		return true, nil
	}
	os.Rename(repoDir+"/.dvs.yaml", repoDir+"/dvs.yaml")

	return true, nil
}
