package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xenciscbc/mysd/internal/config"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Manage docs_to_update list",
	RunE:  runDocsList,
}

var docsAddCmd = &cobra.Command{
	Use:   "add <path>",
	Short: "Add a file path to docs_to_update",
	Args:  cobra.ExactArgs(1),
	RunE:  runDocsAdd,
}

var docsRemoveCmd = &cobra.Command{
	Use:   "remove <path>",
	Short: "Remove a file path from docs_to_update",
	Args:  cobra.ExactArgs(1),
	RunE:  runDocsRemove,
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.AddCommand(docsAddCmd)
	docsCmd.AddCommand(docsRemoveCmd)
}

func runDocsList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if len(cfg.DocsToUpdate) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No docs_to_update configured. Use `mysd docs add <path>` to add files.")
		return nil
	}

	for _, path := range cfg.DocsToUpdate {
		fmt.Fprintln(cmd.OutOrStdout(), path)
	}
	return nil
}

func runDocsAdd(cmd *cobra.Command, args []string) error {
	path := args[0]

	cfg, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Check for duplicate
	for _, existing := range cfg.DocsToUpdate {
		if existing == path {
			fmt.Fprintf(cmd.OutOrStdout(), "%s already configured\n", path)
			return nil
		}
	}

	cfg.DocsToUpdate = append(cfg.DocsToUpdate, path)

	if err := writeDocsToUpdate(cfg.DocsToUpdate); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Added: %s\n", path)
	return nil
}

func runDocsRemove(cmd *cobra.Command, args []string) error {
	path := args[0]

	cfg, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Find and remove path
	newList := make([]string, 0, len(cfg.DocsToUpdate))
	found := false
	for _, existing := range cfg.DocsToUpdate {
		if existing == path {
			found = true
			continue
		}
		newList = append(newList, existing)
	}

	if !found {
		return fmt.Errorf("path not in docs_to_update: %s", path)
	}

	if err := writeDocsToUpdate(newList); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Removed: %s\n", path)
	return nil
}

// writeDocsToUpdate writes the docs_to_update list back to .claude/mysd.yaml,
// preserving all other config fields (same pattern as cmd/model.go runModelSet).
func writeDocsToUpdate(docs []string) error {
	configPath := filepath.Join(".", ".claude", "mysd.yaml")

	v := viper.New()
	v.SetConfigFile(configPath)

	// CRITICAL: ReadInConfig first to preserve existing fields
	if err := v.ReadInConfig(); err != nil {
		// If file doesn't exist, that's OK — SafeWriteConfig will create it
	}

	v.Set("docs_to_update", docs)

	if err := v.WriteConfig(); err != nil {
		// File may not exist yet — try SafeWriteConfig
		if err2 := v.SafeWriteConfig(); err2 != nil {
			return fmt.Errorf("write config: %w", err2)
		}
	}
	return nil
}
