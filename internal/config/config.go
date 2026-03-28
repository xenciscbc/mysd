package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// DefaultModelMap defines the model mapping per profile per agent role.
// Profiles: "quality" | "balanced" | "budget"
// Roles: "spec-writer" | "designer" | "planner" | "executor" | "verifier" | "fast-forward"
//
// Values are short model names ("sonnet", "opus", "haiku") compatible with
// Claude Code's Agent/Task tool model parameter.
var DefaultModelMap = map[string]map[string]string{
	"quality": {
		"spec-writer":     "opus",
		"designer":        "opus",
		"planner":         "opus",
		"executor":        "sonnet",
		"verifier":        "opus",
		"fast-forward":    "sonnet",
		"researcher":      "opus",
		"advisor":         "opus",
		"proposal-writer": "opus",
		"plan-checker":    "opus",
		"reviewer":        "opus",
	},
	"balanced": {
		"spec-writer":     "opus",
		"designer":        "opus",
		"planner":         "opus",
		"executor":        "sonnet",
		"verifier":        "opus",
		"fast-forward":    "sonnet",
		"researcher":      "sonnet",
		"advisor":         "opus",
		"proposal-writer": "sonnet",
		"plan-checker":    "opus",
		"reviewer":        "sonnet",
	},
	"budget": {
		"spec-writer":     "sonnet",
		"designer":        "haiku",
		"planner":         "sonnet",
		"executor":        "haiku",
		"verifier":        "sonnet",
		"fast-forward":    "haiku",
		"researcher":      "sonnet",
		"advisor":         "sonnet",
		"proposal-writer": "sonnet",
		"plan-checker":    "sonnet",
		"reviewer":        "sonnet",
	},
}

// ResolveModel returns the short model name for the given agent role and profile.
// Checks overrides first, then DefaultModelMap[profile][agentRole],
// falling back to "sonnet" if not found.
func ResolveModel(agentRole string, profile string, overrides map[string]string) string {
	if overrides != nil {
		if model, ok := overrides[agentRole]; ok {
			return model
		}
	}
	if profileMap, ok := DefaultModelMap[profile]; ok {
		if model, ok := profileMap[agentRole]; ok {
			return model
		}
	}
	return "sonnet"
}

// Load reads the project configuration from .claude/mysd.yaml (project-level)
// or ~/.claude/mysd.yaml (user-level). If no config file is found, defaults are returned
// without error (convention over config).
func Load(projectRoot string) (ProjectConfig, error) {
	v := viper.New()
	v.SetConfigName("mysd")
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(projectRoot, ".claude"))

	homeDir, err := os.UserHomeDir()
	if err == nil {
		v.AddConfigPath(filepath.Join(homeDir, ".claude"))
	}

	// Set defaults from Defaults()
	d := Defaults()
	v.SetDefault("execution_mode", d.ExecutionMode)
	v.SetDefault("agent_count", d.AgentCount)
	v.SetDefault("atomic_commits", d.AtomicCommits)
	v.SetDefault("tdd", d.TDD)
	v.SetDefault("test_generation", d.TestGeneration)
	v.SetDefault("response_language", d.ResponseLanguage)
	v.SetDefault("document_language", d.DocumentLanguage)
	v.SetDefault("model_profile", d.ModelProfile)
	v.SetDefault("worktree_dir", d.WorktreeDir)
	v.SetDefault("auto_mode", d.AutoMode)
	v.SetDefault("docs_to_update", d.DocsToUpdate)

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return ProjectConfig{}, err
		}
		// Config file not found — convention over config, use defaults
	}

	var cfg ProjectConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return ProjectConfig{}, err
	}
	return cfg, nil
}

// BindFlags binds cobra/pflag flags to viper config keys.
// This is called by cmd/root.go to allow flag override of config values.
func BindFlags(v *viper.Viper, flags *pflag.FlagSet) {
	_ = v.BindPFlag("execution_mode", flags.Lookup("execution-mode"))
	_ = v.BindPFlag("agent_count", flags.Lookup("agent-count"))
	_ = v.BindPFlag("atomic_commits", flags.Lookup("atomic-commits"))
	_ = v.BindPFlag("tdd", flags.Lookup("tdd"))
	_ = v.BindPFlag("response_language", flags.Lookup("lang"))
}
