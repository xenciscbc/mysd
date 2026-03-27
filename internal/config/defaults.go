package config

// ProjectConfig holds all user-configurable preferences for my-ssd.
type ProjectConfig struct {
	ExecutionMode    string            `yaml:"execution_mode" mapstructure:"execution_mode"`       // "single" | "wave"
	AgentCount       int               `yaml:"agent_count" mapstructure:"agent_count"`             // default 1
	AtomicCommits    bool              `yaml:"atomic_commits" mapstructure:"atomic_commits"`       // default false
	TDD              bool              `yaml:"tdd" mapstructure:"tdd"`                             // default false
	TestGeneration   bool              `yaml:"test_generation" mapstructure:"test_generation"`     // default false
	ResponseLanguage string            `yaml:"response_language" mapstructure:"response_language"` // e.g. "zh-TW"
	DocumentLanguage string            `yaml:"document_language" mapstructure:"document_language"` // e.g. "en"
	ModelProfile     string            `yaml:"model_profile" mapstructure:"model_profile"`         // "quality" | "balanced" | "budget"
	ModelOverrides   map[string]string `yaml:"model_overrides" mapstructure:"model_overrides"`     // per-agent model overrides
	WorktreeDir      string            `yaml:"worktree_dir" mapstructure:"worktree_dir"`           // default ".worktrees"
	AutoMode         bool              `yaml:"auto_mode" mapstructure:"auto_mode"`                 // default false
	DocsToUpdate     []string          `yaml:"docs_to_update" mapstructure:"docs_to_update"`       // files to update after archive
	StatuslineEnabled *bool            `yaml:"statusline_enabled,omitempty" mapstructure:"statusline_enabled"` // nil = not set (treated as true)
}

// Defaults returns a ProjectConfig with convention-over-config default values.
func Defaults() ProjectConfig {
	return ProjectConfig{
		ExecutionMode:    "single",
		AgentCount:       1,
		AtomicCommits:    false,
		TDD:              false,
		TestGeneration:   false,
		ResponseLanguage: "",
		DocumentLanguage: "",
		ModelProfile:     "balanced",
		ModelOverrides:   nil,
		WorktreeDir:      ".worktrees",
		AutoMode:         false,
		DocsToUpdate:     nil,
	}
}
