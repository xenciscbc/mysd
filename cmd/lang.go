package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xenciscbc/mysd/internal/config"
	"github.com/xenciscbc/mysd/internal/spec"
)

var langCmd = &cobra.Command{
	Use:   "lang",
	Short: "Display or set language settings",
	RunE:  runLangRead,
}

var langSetCmd = &cobra.Command{
	Use:   "set <locale>",
	Short: "Set response language and OpenSpec locale (BCP47 format, e.g. zh-TW, en-US)",
	Args:  cobra.ExactArgs(1),
	RunE:  runLangSet,
}

func init() {
	langCmd.AddCommand(langSetCmd)
	rootCmd.AddCommand(langCmd)
}

func runLangRead(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	osCfg, err := spec.ReadOpenSpecConfig(".")
	if err != nil {
		return fmt.Errorf("read openspec config: %w", err)
	}

	responseLang := cfg.ResponseLanguage
	if responseLang == "" {
		responseLang = "(not set)"
	}
	locale := osCfg.Locale
	if locale == "" {
		locale = "(not set)"
	}

	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "Language settings:\n")
	fmt.Fprintf(w, "  mysd.yaml response_language: %s\n", responseLang)
	fmt.Fprintf(w, "  openspec/config.yaml locale: %s\n", locale)

	return nil
}

func runLangSet(cmd *cobra.Command, args []string) error {
	locale := args[0]

	// Step 1: Read current values for potential rollback
	cfg, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	oldResponseLang := cfg.ResponseLanguage

	osCfg, err := spec.ReadOpenSpecConfig(".")
	if err != nil {
		return fmt.Errorf("read openspec config: %w", err)
	}

	// Step 2: Write mysd.yaml first (via Viper), preserving other fields (Pitfall 1)
	configPath := filepath.Join(".", ".claude", "mysd.yaml")
	v := viper.New()
	v.SetConfigFile(configPath)
	_ = v.ReadInConfig() // preserve existing fields; ignore error if file absent

	v.Set("response_language", locale)

	if err := v.WriteConfig(); err != nil {
		// File may not exist yet — try SafeWriteConfig
		if err2 := v.SafeWriteConfig(); err2 != nil {
			return fmt.Errorf("write mysd.yaml: %w", err2)
		}
	}

	// Step 3: Write openspec/config.yaml; rollback mysd.yaml on failure
	osCfg.Locale = locale
	if err := spec.WriteOpenSpecConfig(".", osCfg); err != nil {
		// ROLLBACK: restore mysd.yaml to old value
		v.Set("response_language", oldResponseLang)
		_ = v.WriteConfig()
		return fmt.Errorf("write openspec config (rolled back mysd.yaml): %w", err)
	}

	// Step 4: Print success confirmation
	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "Language updated:\n")
	fmt.Fprintf(w, "  response_language: %s\n", locale)
	fmt.Fprintf(w, "  locale: %s\n", locale)

	return nil
}
