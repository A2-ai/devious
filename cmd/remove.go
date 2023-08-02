package cmd

import (
	"github.com/spf13/cobra"
)

func runRemoveCmd(cmd *cobra.Command, args []string) error {
	return nil
}

func getRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a file from storage and delete its metadata",
		RunE:  runRemoveCmd,
	}

	return cmd
}
