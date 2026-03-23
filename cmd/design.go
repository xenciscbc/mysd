package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var designCmd = &cobra.Command{
	Use:   "design",
	Short: "Capture technical decisions and architecture",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(cmd.OutOrStdout(), "not yet implemented")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(designCmd)
}
