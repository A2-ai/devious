package cmd

import (
	"dvs/internal/log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var verbose bool

func getRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dvs",
		Short: "👺 Devious\nA file linker that enables you to work with large files while keeping them under version control.\nSee https://github.com/A2-ai/devious for more information.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Configure logger
			var logLevel slog.Level
			if verbose {
				logLevel = slog.LevelDebug
			} else {
				logLevel = slog.LevelInfo
			}
			log.ConfigureGlobalLogger(logLevel)

			return nil
		},
	}

	// Add global flags
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")

	// Add commands
	cmd.AddCommand(getInitCmd())
	cmd.AddCommand(getStatusCmd())
	cmd.AddCommand(getAddCmd())
	cmd.AddCommand(getGetCmd())
	cmd.AddCommand(getRemoveCmd())

	return cmd
}

func Execute() {
	err := getRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
