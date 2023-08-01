package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runGetCmd(cmd *cobra.Command, args []string) error {
	return nil
}

func getGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a file from storage",
		RunE:  runGetCmd,
	}

	cmd.Flags().BoolP("force", "f", false, "Force remove a file from the index")
	viper.BindPFlag("force", cmd.Flags().Lookup("force"))

	return cmd
}
