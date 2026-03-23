package cmd

import (
	"os"
	"path/filepath"

	"github.com/mysd/internal/config"
	"github.com/mysd/internal/output"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize project config (.claude/mysd.yaml)",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().Bool("force", false, "overwrite existing config")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	p := output.NewPrinter(cmd.OutOrStdout())

	force, _ := cmd.Flags().GetBool("force")
	targetPath := filepath.Join(".", ".claude", "mysd.yaml")

	// Check if file already exists
	if _, err := os.Stat(targetPath); err == nil && !force {
		p.Warning("Config already exists: .claude/mysd.yaml (use --force to overwrite)")
		return nil
	}

	// Ensure .claude/ directory exists
	if err := os.MkdirAll(filepath.Join(".", ".claude"), 0755); err != nil {
		p.Error("Failed to create .claude/ directory: " + err.Error())
		return err
	}

	// Get defaults and marshal to YAML
	cfg := config.Defaults()
	data, err := yaml.Marshal(cfg)
	if err != nil {
		p.Error("Failed to marshal config: " + err.Error())
		return err
	}

	// Prepend descriptive comment header
	header := "# my-ssd project configuration\n# See: https://github.com/mysd/docs/config\n"
	content := []byte(header + string(data))

	if err := os.WriteFile(targetPath, content, 0644); err != nil {
		p.Error("Failed to write config: " + err.Error())
		return err
	}

	p.Success("Created config: .claude/mysd.yaml")
	return nil
}
