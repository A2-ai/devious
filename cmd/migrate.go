package cmd

import (
	"dvs/internal/log"

	"github.com/spf13/cobra"
)

func runMigrateCmd(cmd *cobra.Command, args []string) {
	log.Print("Migrating storage...")

	// err := migrate()
	// if err != nil {
	// 	log.Error("Failed to migrate storage", err)
	// 	return
	// }

	// log.Info("Successfully migrated storage")
}

func getMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrates storage to the latest format",
		Long:  "Gets the status of files tracked by devious. If glob(s) are provided, only those globs will be checked. Otherwise, all files in the current git repository will be checked.",
		PreRun: func(cmd *cobra.Command, args []string) {
			log.PrintLogo()
		},
		Run: runMigrateCmd,
	}

	return cmd
}
