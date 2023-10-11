package migrate

func MigrateToLatest() {
	MigrateMetaFiles()
	MigrateStorageFiles()
}
