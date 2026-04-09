package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// DefaultModelMap defines the model mapping per profile per agent role.
// Profiles: "quality" | "balanced" | "budget"
// Roles: "spec-writer" | "designer" | "planner" | "executor" | "spec-executor" | "verifier" | "fast-forward"
//
// Values are short model names ("sonnet", "opus", "haiku") compatible with
// Claude Code's Agent/Task tool model parameter.
var DefaultModelMap = map[string]map[string]string{
	"quality": {
		"spec-writer":     "opus",
		"designer":        "opus",
		"planner":         "opus",
		"executor":        "sonnet",
		"spec-executor":   "opus",
		"verifier":        "opus",
		"fast-forward":    "sonnet",
		"researcher":      "opus",
		"advisor":         "opus",
		"proposal-writer": "opus",
		"plan-checker":    "opus",
		"reviewer":        "opus",
		"scanner":         "opus",
		"uat-guide":       "opus",
	},
	"balanced": {
		"spec-writer":     "opus",
		"designer":        "opus",
		"planner":         "opus",
		"executor":        "sonnet",
		"spec-executor":   "opus",
		"verifier":        "opus",
		"fast-forward":    "sonnet",
		"researcher":      "sonnet",
		"advisor":         "opus",
		"proposal-writer": "sonnet",
		"plan-checker":    "opus",
		"reviewer":        "sonnet",
		"scanner":         "sonnet",
		"uat-guide":       "sonnet",
	},
	"budget": {
		"spec-writer":     "sonnet",
		"designer":        "haiku",
		"planner":         "sonnet",
		"executor":        "haiku",
		"spec-executor":   "sonnet",
		"verifier":        "sonnet",
		"fast-forward":    "haiku",
		"researcher":      "sonnet",
		"advisor":         "sonnet",
		"proposal-writer": "sonnet",
		"plan-checker":    "sonnet",
		"reviewer":        "sonnet",
		"scanner":         "sonnet",
		"uat-guide":       "sonnet",
	},
}

// ResolveModel returns the short model name for the given agent role and profile.
// Resolution order:
//  1. ModelOverrides[role]
//  2. CustomProfiles[profile].Models[role]
//  3. DefaultModelMap[CustomProfiles[profile].Base][role]
//  4. DefaultModelMap[profile][role] (when profile is a built-in)
//  5. "sonnet" (fallback)
func ResolveModel(agentRole string, profile string, overrides map[string]string, customProfiles map[string]CustomProfile) string {
	if overrides != nil {
		if model, ok := overrides[agentRole]; ok {
			return model
		}
	}
	if customProfiles != nil {
		if cp, ok := customProfiles[profile]; ok {
			if model, ok := cp.Models[agentRole]; ok {
				return model
			}
			if baseMap, ok := DefaultModelMap[cp.Base]; ok {
				if model, ok := baseMap[agentRole]; ok {
					return model
				}
			}
			return "sonnet"
		}
	}
	if profileMap, ok := DefaultModelMap[profile]; ok {
		if model, ok := profileMap[agentRole]; ok {
			return model
		}
	}
	return "sonnet"
}

// ValidateCustomProfiles checks custom profiles for unknown role names and invalid base profiles.
// Returns a list of warning messages (empty if no issues found).
func ValidateCustomProfiles(knownRoles []string, customProfiles map[string]CustomProfile) []string {
	var warnings []string
	roleSet := make(map[string]bool, len(knownRoles))
	for _, r := range knownRoles {
		roleSet[r] = true
	}
	for name, cp := range customProfiles {
		if cp.Base != "" {
			if _, ok := DefaultModelMap[cp.Base]; !ok {
				warnings = append(warnings, fmt.Sprintf("custom profile %q: base %q is not a valid built-in profile", name, cp.Base))
			}
		}
		for role := range cp.Models {
			if !roleSet[role] {
				warnings = append(warnings, fmt.Sprintf("custom profile %q: role %q is not a known role", name, role))
			}
		}
	}
	return warnings
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
