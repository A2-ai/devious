package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dvs",
	Short: "Devious - a CLI for dealing with large files",
}

func getRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dvs",
		Short: "Devious - a CLI for dealing with large files",
	}

	cmd.AddCommand(getAddCmd())
	cmd.AddCommand(getRemoveCmd())
	cmd.AddCommand(getListCmd())
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
