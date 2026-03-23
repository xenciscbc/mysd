package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Run tasks with pre-execution alignment",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(cmd.OutOrStdout(), "not yet implemented")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
}
