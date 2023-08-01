package cmd

import (
	"github.com/spf13/cobra"
)

func runListCmd(cmd *cobra.Command, args []string) error {
	return nil
}

func getListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists files under the influence of devious",
		RunE:  runListCmd,
	}

	return cmd
}
