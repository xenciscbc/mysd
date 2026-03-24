package cmd

import (
	"fmt"

	"github.com/xenciscbc/mysd/internal/output"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/spf13/cobra"
)

var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Extract changes from current conversation into propose flow",
	RunE:  runCapture,
}

func init() {
	captureCmd.Flags().String("name", "", "pre-set change name and scaffold the directory")
	rootCmd.AddCommand(captureCmd)
}

func runCapture(cmd *cobra.Command, args []string) error {
	p := output.NewPrinter(cmd.OutOrStdout())

	name, _ := cmd.Flags().GetString("name")

	if name != "" {
		// Scaffold the change directory so SKILL.md can proceed immediately
		specDir, _, err := spec.DetectSpecDir(".")
		if err != nil {
			specDir = ".specs"
		}

		change, err := spec.Scaffold(name, specDir)
		if err != nil {
			return fmt.Errorf("scaffold change: %w", err)
		}

		p.Success(fmt.Sprintf("Created change directory: %s", change.Dir))
		p.Info("Use /mysd:capture in Claude Code to extract changes from conversation")
		return nil
	}

	// No --name: print guidance for SKILL.md layer
	// Actual conversation analysis requires Claude Code context (Pitfall 6)
	p.Info("Use /mysd:capture in Claude Code to extract changes from conversation")
	p.Muted("Optional: mysd capture --name <change-name> to pre-scaffold the directory")
	return nil
}
