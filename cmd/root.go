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
		Short: "Devious - a CLI for dealing with large files",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Configure logger
			var logLevel slog.Level
			if verbose {
				logLevel = slog.LevelDebug
			} else {
				logLevel = slog.LevelInfo
			}
			log.SetGlobalLogLevel(logLevel)

			return nil
		},
	}

	// Add global flags
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")

	// Add commands
	cmd.AddCommand(getAddCmd())
	cmd.AddCommand(getRemoveCmd())
	cmd.AddCommand(getStatusCmd())
	cmd.AddCommand(getGetCmd())
	cmd.AddCommand(getInitCmd())

	return cmd
}

func Execute() {
	err := getRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
