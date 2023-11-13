package migrate

func MigrateToLatest(dry bool) (match bool, err error) {
	files, err := MigrateMetaFiles(dry)
	if err != nil {
		return false, err
	}
	if len(files) > 0 {
		return true, err
	}

	files, err = MigrateStorageFiles(dry)
	if err != nil {
		return false, err
	}
	if len(files) > 0 {
		return true, err
	}

	return false, nil
}
