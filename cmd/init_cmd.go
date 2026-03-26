package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xenciscbc/mysd/internal/config"
	"github.com/xenciscbc/mysd/internal/output"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize openspec structure and project config",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	p := output.NewPrinter(cmd.OutOrStdout())

	// Scaffold openspec/ and openspec/specs/ directories (idempotent via MkdirAll)
	if err := scaffoldOpenSpecDir("."); err != nil {
		p.Error("Failed to scaffold openspec structure: " + err.Error())
		return err
	}

	// Ensure .claude/ directory exists
	if err := os.MkdirAll(filepath.Join(".", ".claude"), 0755); err != nil {
		p.Error("Failed to create .claude/ directory: " + err.Error())
		return err
	}

	// Create default .claude/mysd.yaml if not present
	targetPath := filepath.Join(".", ".claude", "mysd.yaml")
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		cfg := config.Defaults()
		data, marshalErr := yaml.Marshal(cfg)
		if marshalErr != nil {
			p.Error("Failed to marshal config: " + marshalErr.Error())
			return marshalErr
		}
		header := "# my-ssd project configuration\n# See: https://github.com/mysd/docs/config\n"
		content := []byte(header + string(data))
		if writeErr := os.WriteFile(targetPath, content, 0644); writeErr != nil {
			p.Error("Failed to write config: " + writeErr.Error())
			return writeErr
		}
	}

	p.Success("Initialized openspec structure. Run /mysd:scan to discover codebase.")
	return nil
}
