package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xenciscbc/mysd/internal/config"
	"github.com/xenciscbc/mysd/internal/output"
)

// knownRoles lists all agent roles in deterministic order for consistent output.
var knownRoles = []string{
	"spec-writer", "designer", "planner", "executor", "verifier",
	"fast-forward", "researcher", "advisor", "proposal-writer", "plan-checker",
}

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Display or set model profile",
	RunE:  runModelRead,
}

var modelSetCmd = &cobra.Command{
	Use:   "set <profile>",
	Short: "Set model profile (quality, balanced, budget)",
	Args:  cobra.ExactArgs(1),
	RunE:  runModelSet,
}

func init() {
	modelCmd.AddCommand(modelSetCmd)
	rootCmd.AddCommand(modelCmd)
}

func runModelRead(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	profile := cfg.ModelProfile
	if profile == "" {
		profile = "balanced"
	}

	w := cmd.OutOrStdout()
	_ = output.NewPrinter(w) // used for TTY detection only; output via fmt.Fprintf for plain text

	fmt.Fprintf(w, "Profile: %s\n\n", profile)
	fmt.Fprintf(w, "%-20s %s\n", "Role", "Model")
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 40))
	for _, role := range knownRoles {
		model := config.ResolveModel(role, profile, cfg.ModelOverrides)
		fmt.Fprintf(w, "%-20s %s\n", role, model)
	}

	return nil
}

func runModelSet(cmd *cobra.Command, args []string) error {
	profile := args[0]

	// Validate profile against DefaultModelMap
	if _, ok := config.DefaultModelMap[profile]; !ok {
		return fmt.Errorf("unknown profile %q; valid profiles: quality, balanced, budget", profile)
	}

	configPath := filepath.Join(".", ".claude", "mysd.yaml")

	v := viper.New()
	v.SetConfigFile(configPath)

	// CRITICAL: ReadInConfig first to preserve existing fields (Pitfall 1)
	if err := v.ReadInConfig(); err != nil {
		// If file doesn't exist, that's OK — SafeWriteConfig will create it
	}

	v.Set("model_profile", profile)

	if err := v.WriteConfig(); err != nil {
		// File may not exist yet — try SafeWriteConfig
		if err2 := v.SafeWriteConfig(); err2 != nil {
			return fmt.Errorf("write config: %w", err2)
		}
	}

	p := output.NewPrinter(cmd.OutOrStdout())
	p.Success(fmt.Sprintf("Model profile set to: %s", profile))

	return nil
}
