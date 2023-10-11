package cmd

import (
	"dvs/internal/log"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"

func getRootCmd() *cobra.Command {
	var json bool

	cmd := &cobra.Command{
		Use:   "dvs",
		Short: "🌀 Devious\nA file linker that allows you to version large files under Git.\nSee https://github.com/A2-ai/devious for more information.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if json {
				log.JsonLogging = true
			}
		},
	}

	// Version
	cmd.Version = Version
	cmd.SetVersionTemplate("{{.Version}}\n")

	// Disable completion command
	cmd.CompletionOptions.DisableDefaultCmd = true

	// Add global flags
	cmd.PersistentFlags().BoolVarP(&json, "json", "j", false, "log in JSON format")

	// Add commands
	cmd.AddCommand(getInitCmd())
	cmd.AddCommand(getStatusCmd())
	cmd.AddCommand(getAddCmd())
	cmd.AddCommand(getGetCmd())
	cmd.AddCommand(getRemoveCmd())
	cmd.AddCommand(getMigrateCmd())

	return cmd
}

func Execute() {
	err := getRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
