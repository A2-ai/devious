package cmd

import (
	"dvs/internal/log"
	"dvs/internal/migrate"
	"fmt"

	"github.com/spf13/cobra"
)

func runMigrateModule(msg string, migrateFunc func(bool) ([]string, error), dry bool) {
	log.Print(msg)
	files, err := migrateFunc(dry)
	log.OverwritePreviousLine()
	if err != nil {
		log.Print(msg, log.IconFailure)
		log.Print("  ", log.ColorRed(err))
	} else if len(files) == 0 {
		log.Print(msg, log.IconSuccess, log.ColorGreen("already up to date"))
	} else {
		log.Print(msg, log.IconSuccess, log.ColorGreen("migrated ", fmt.Sprint(len(files)), " files"))
		for _, file := range files {
			log.Print("  ", log.ColorFile(file))
		}
	}
}

func runMigrateCmd(cmd *cobra.Command, args []string) {
	dry, err := cmd.Flags().GetBool("dry")
	if err != nil {
		return
	}

	runMigrateModule("Migrating local metadata...", migrate.MigrateMetaFiles, dry)
	runMigrateModule("Migrating files in storage...", migrate.MigrateStorageFiles, dry)

	log.Print("\nMigration complete!")

	log.DumpAndExit(0)
}

func getMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "migrate",
		Short:            "Migrates data to the latest format",
		Long:             "Migrates data to the latest format. This includes the meta files and the storage.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		PreRun: func(cmd *cobra.Command, args []string) {
			log.PrintLogo()
		},
		Run: runMigrateCmd,
	}

	// Add --dry flag
	cmd.Flags().BoolP("dry", "d", false, "dry run")

	return cmd
}
