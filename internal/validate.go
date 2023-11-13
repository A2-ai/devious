package internal

import (
	"dvs/internal/git"
	"dvs/internal/log"
	"dvs/internal/migrate"
)

func ValidateEnvironment() (bool, string) {
	// Ensure we can find a git repository
	_, err := git.GetNearestRepoDir(".")
	if err != nil {
		return false, log.IconFailure + "Failed to find a git repository here -- are you sure you're in one?"
	}

	// Ensure data is migrated
	match, err := migrate.MigrateToLatest(true)
	if match || err != nil {
		return false, log.IconWarning + " Data is not migrated - run " + log.ColorFaint("dvs migrate") + " to migrate your data to the current version"
	}

	return true, ""
}
