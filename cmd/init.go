package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runInitCmd(cmd *cobra.Command, args []string) error {
	return nil
}

func getInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize devious",
		RunE:  runInitCmd,
	}

	cmd.Flags().BoolP("force", "f", false, "Force remove a file from the index")
	viper.BindPFlag("force", cmd.Flags().Lookup("force"))

	return cmd
}
