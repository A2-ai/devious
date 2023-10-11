package cmd

import (
	"dvs/internal/log"
	"dvs/internal/migrate"
	"fmt"

	"github.com/spf13/cobra"
)

func runMigrateModule(msg string, migrateFunc func() ([]string, error)) {
	log.Print(msg)
	filesModified, err := migrateFunc()
	log.OverwritePreviousLine()
	if err != nil {
		log.Print(msg, log.IconFailure)
		log.Print("  ", log.ColorRed(err))
	} else if len(filesModified) == 0 {
		log.Print(msg, log.IconSuccess, log.ColorGreen("already up to date"))
	} else {
		log.Print(msg, log.IconSuccess, log.ColorGreen("migrated ", fmt.Sprint(len(filesModified)), " files"))
		for _, file := range filesModified {
			log.Print("  ", log.ColorFile(file))
		}
	}
}

func runMigrateCmd(cmd *cobra.Command, args []string) {
	runMigrateModule("Migrating local metadata...", migrate.MigrateMetaFiles)
	runMigrateModule("Migrating files in storage...", migrate.MigrateStorageFiles)

	log.Print("\nMigration complete!")
}

func getMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrates data to the latest format",
		Long:  "Migrates data to the latest format. This includes the meta files and the storage.",
		PreRun: func(cmd *cobra.Command, args []string) {
			log.PrintLogo()
		},
		Run: runMigrateCmd,
	}

	return cmd
}
