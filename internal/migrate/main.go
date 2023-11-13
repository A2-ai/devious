package migrate

func MigrateToLatest(dry bool) (match bool, err error) {
	match = false

	files, err := MigrateConfig(dry)
	if err != nil {
		return false, err
	}
	if len(files) > 0 {
		match = true
	}

	files, err = MigrateMetaFiles(dry)
	if err != nil {
		return false, err
	}
	if len(files) > 0 {
		match = true
	}

	files, err = MigrateStorageFiles(dry)
	if err != nil {
		return false, err
	}
	if len(files) > 0 {
		match = true
	}

	return match, nil
}
