package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runAddCmd(cmd *cobra.Command, args []string) error {
	return nil
}

func getAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a file to storage",
		RunE:  runRemoveCmd,
	}

	cmd.Flags().BoolP("force", "f", false, "Force remove a file from the index")
	viper.BindPFlag("force", cmd.Flags().Lookup("force"))

	return cmd
}
