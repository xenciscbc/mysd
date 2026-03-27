package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statuslineCmd = &cobra.Command{
	Use:   "statusline [on|off]",
	Short: "Toggle or set statusline display",
	RunE:  runStatusline,
}

func init() {
	rootCmd.AddCommand(statuslineCmd)
}

// runStatusline is the cobra RunE handler for the statusline command.
// It delegates to runStatuslineInDir using the current directory.
func runStatusline(cmd *cobra.Command, args []string) error {
	return runStatuslineInDir(cmd, args, ".")
}

// runStatuslineInDir implements statusline on/off/toggle using the given base directory.
// Separated for testability — tests pass t.TempDir() as baseDir.
func runStatuslineInDir(cmd *cobra.Command, args []string, baseDir string) error {
	configPath := filepath.Join(baseDir, ".claude", "mysd.yaml")

	v := viper.New()
	v.SetConfigFile(configPath)
	_ = v.ReadInConfig()

	var newValue bool

	if len(args) == 0 {
		// Toggle: if not set, default is true (D-12), toggle -> false; if set, flip.
		if !v.IsSet("statusline_enabled") {
			newValue = false
		} else {
			newValue = !v.GetBool("statusline_enabled")
		}
	} else {
		switch args[0] {
		case "on":
			newValue = true
		case "off":
			newValue = false
		default:
			return fmt.Errorf("invalid argument %q; use: on, off, or no argument to toggle", args[0])
		}
	}

	v.Set("statusline_enabled", newValue)

	// Ensure the .claude/ directory exists before writing config.
	claudeDir := filepath.Join(baseDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	if err := v.WriteConfig(); err != nil {
		if err2 := v.SafeWriteConfig(); err2 != nil {
			return fmt.Errorf("write config: %w", err2)
		}
	}

	onOff := "off"
	if newValue {
		onOff = "on"
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Statusline: %s\n", onOff)

	return nil
}
