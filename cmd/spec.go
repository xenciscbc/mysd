package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Define detailed requirements with RFC 2119 keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(cmd.OutOrStdout(), "not yet implemented")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(specCmd)
}
