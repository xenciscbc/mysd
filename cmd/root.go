package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "mysd",
	Short: "Spec-Driven Development for AI programming",
	Long: `my-ssd integrates OpenSpec's Spec-Driven Development methodology
with a GSD-level planning/execution/verification engine.

It enables structured spec-driven AI programming for solo developers.`,
}

var cfgFile string

// Execute runs the root command. Called by main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default .claude/mysd.yaml)")
	rootCmd.PersistentFlags().String("execution-mode", "", "execution mode: single or wave")
	rootCmd.PersistentFlags().Int("agent-count", 0, "number of agents for wave mode")
	rootCmd.PersistentFlags().String("lang", "", "response language override")
	rootCmd.PersistentFlags().String("doc-lang", "", "document output language override")
	rootCmd.PersistentFlags().Bool("tdd", false, "enable TDD mode")
	rootCmd.PersistentFlags().Bool("atomic-commits", false, "enable atomic git commits per task")

	_ = viper.BindPFlag("execution_mode", rootCmd.PersistentFlags().Lookup("execution-mode"))
	_ = viper.BindPFlag("agent_count", rootCmd.PersistentFlags().Lookup("agent-count"))
	_ = viper.BindPFlag("response_language", rootCmd.PersistentFlags().Lookup("lang"))
	_ = viper.BindPFlag("document_language", rootCmd.PersistentFlags().Lookup("doc-lang"))
	_ = viper.BindPFlag("tdd", rootCmd.PersistentFlags().Lookup("tdd"))
	_ = viper.BindPFlag("atomic_commits", rootCmd.PersistentFlags().Lookup("atomic-commits"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("mysd")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".claude")

		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(filepath.Join(home, ".claude"))
		}
	}

	viper.AutomaticEnv()

	// Ignore errors (including ConfigFileNotFoundError) — convention over config
	_ = viper.ReadInConfig()
}
