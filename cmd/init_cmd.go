package cmd

import (
	"encoding/json"
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

	// Install statusline hook from embedded content (D-06, D-08)
	hookDir := filepath.Join(".", ".claude", "hooks")
	if err := os.MkdirAll(hookDir, 0755); err != nil {
		p.Error("Failed to create hooks directory: " + err.Error())
		return err
	}
	hookDest := filepath.Join(hookDir, "mysd-statusline.js")
	if err := os.WriteFile(hookDest, statuslineHookBytes, 0644); err != nil {
		p.Error("Failed to install statusline hook: " + err.Error())
		return err
	}

	// Write statusLine key to .claude/settings.json (D-06, D-07)
	if err := writeSettingsStatusLine(filepath.Join(".", ".claude")); err != nil {
		p.Error("Failed to write settings.json: " + err.Error())
		return err
	}

	p.Success("Initialized openspec structure. Statusline configured. Run /mysd:scan to discover codebase.")
	return nil
}

// writeSettingsStatusLine reads (or creates) .claude/settings.json, sets the
// statusLine key, and writes it back, preserving all other existing keys (D-07).
func writeSettingsStatusLine(claudeDir string) error {
	settingsPath := filepath.Join(claudeDir, "settings.json")

	// Read existing (or start fresh) — D-07: merge, preserve other keys
	raw := map[string]interface{}{}
	if data, err := os.ReadFile(settingsPath); err == nil {
		_ = json.Unmarshal(data, &raw) // silent fail on parse error
	}

	// Set only the statusLine key
	raw["statusLine"] = map[string]interface{}{
		"type":    "command",
		"command": "node .claude/hooks/mysd-statusline.js",
	}

	out, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(settingsPath, append(out, '\n'), 0644)
}
