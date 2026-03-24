package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/xenciscbc/mysd/internal/scanner"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan codebase and output metadata for spec generation",
	RunE:  runScan,
}

func init() {
	scanCmd.Flags().Bool("context-only", false, "Output scan context as JSON (for /mysd:scan agent consumption)")
	scanCmd.Flags().StringSlice("exclude", nil, "Directories to exclude from scan")
	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) error {
	contextOnly, _ := cmd.Flags().GetBool("context-only")
	if !contextOnly {
		return fmt.Errorf("usage: mysd scan --context-only [--exclude dir1,dir2]\nDirect execution via /mysd:scan")
	}
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	return runScanContextOnly(cmd.OutOrStdout(), ".", exclude)
}

func runScanContextOnly(out io.Writer, root string, exclude []string) error {
	ctx, err := scanner.BuildScanContext(root, exclude)
	if err != nil {
		return fmt.Errorf("build scan context: %w", err)
	}
	data, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal scan context: %w", err)
	}
	_, err = out.Write(data)
	return err
}
