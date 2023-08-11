package cmd

import (
	"dvs/internal/log"
	"os"

	"github.com/spf13/cobra"
)

var json bool

func getRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dvs",
		Short: "👺 Devious\nA file linker that enables you to work with large files while keeping them under version control.\nSee https://github.com/A2-ai/devious for more information.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if json {
				log.JsonLogger = &log.JsonLog{
					Files: make(map[string]log.JsonFile),
				}
			}
		},
	}

	// Add global flags
	cmd.PersistentFlags().BoolVarP(&json, "json", "j", false, "log in JSON format")

	// Add commands
	cmd.AddCommand(getInitCmd())
	cmd.AddCommand(getStatusCmd())
	cmd.AddCommand(getAddCmd())
	cmd.AddCommand(getGetCmd())
	cmd.AddCommand(getRemoveCmd())

	return cmd
}

func Execute() {
	defer log.Dump()

	err := getRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
