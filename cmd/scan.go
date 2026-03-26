package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xenciscbc/mysd/internal/output"
	"github.com/xenciscbc/mysd/internal/scanner"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan codebase and output metadata for spec generation",
	RunE:  runScan,
}

func init() {
	scanCmd.Flags().Bool("context-only", false, "Output scan context as JSON (for /mysd:scan agent consumption)")
	scanCmd.Flags().Bool("scaffold-only", false, "Create openspec/ directory structure without scanning")
	scanCmd.Flags().StringSlice("exclude", nil, "Directories to exclude from scan")
	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) error {
	contextOnly, _ := cmd.Flags().GetBool("context-only")
	scaffoldOnly, _ := cmd.Flags().GetBool("scaffold-only")

	if scaffoldOnly {
		return runScanScaffoldOnly(cmd)
	}
	if contextOnly {
		exclude, _ := cmd.Flags().GetStringSlice("exclude")
		return runScanContextOnly(cmd.OutOrStdout(), ".", exclude)
	}
	return fmt.Errorf("usage: mysd scan --context-only|--scaffold-only [--exclude dir1,dir2]")
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

func runScanScaffoldOnly(cmd *cobra.Command) error {
	p := output.NewPrinter(cmd.OutOrStdout())
	if err := scaffoldOpenSpecDir("."); err != nil {
		p.Error("Failed to scaffold openspec structure: " + err.Error())
		return err
	}
	p.Success("Initialized openspec structure. Run /mysd:scan to discover codebase.")
	return nil
}

// scaffoldOpenSpecDir creates the openspec/ and openspec/specs/ directories.
// It does NOT create openspec/config.yaml (per D-06, locale is set by SKILL.md agent).
// This function is idempotent — safe to call multiple times via os.MkdirAll.
func scaffoldOpenSpecDir(root string) error {
	dirs := []string{
		filepath.Join(root, "openspec"),
		filepath.Join(root, "openspec", "specs"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
	}
	return nil
}
