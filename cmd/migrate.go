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

		log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
			Severity: "error",
			Message:  err.Error(),
		})
	} else if len(files) == 0 {
		log.Print(msg, log.IconSuccess, log.ColorGreen("already up to date"))

		log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
			Severity: "info",
			Message:  "already up to date",
		})
	} else {
		log.Print(msg, log.IconSuccess, log.ColorGreen("migrated ", fmt.Sprint(len(files)), " files"))
		for _, file := range files {
			log.Print("  ", log.ColorFile(file))
		}

		log.JsonLogger.Issues = append(log.JsonLogger.Issues, log.JsonIssue{
			Severity: "info",
			Message:  "migrated " + fmt.Sprint(len(files)) + " files",
		})
	}
}

func runMigrateCmd(cmd *cobra.Command, args []string) {
	dry, err := cmd.Flags().GetBool("dry")
	if err != nil {
		return
	}

	runMigrateModule("Migrating config...", migrate.MigrateConfig, dry)
	runMigrateModule("Migrating local metadata...", migrate.MigrateMetaFiles, dry)
	runMigrateModule("Migrating files in storage...", migrate.MigrateStorageFiles, dry)

	log.DumpAndExit(0)
}

func getMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrates data to the latest format",
		Long:  "Migrates data to the latest format. This includes the meta files and the storage.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Flags().Changed("json") {
				log.JsonLogging = true
			}
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			log.PrintLogo()
		},
		Run: runMigrateCmd,
	}

	// Add --dry flag
	cmd.Flags().BoolP("dry", "d", false, "dry run")

	return cmd
}
