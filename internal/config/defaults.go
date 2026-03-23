package config

// ProjectConfig holds all user-configurable preferences for my-ssd.
type ProjectConfig struct {
	ExecutionMode    string `yaml:"execution_mode" mapstructure:"execution_mode"`      // "single" | "wave"
	AgentCount       int    `yaml:"agent_count" mapstructure:"agent_count"`            // default 1
	AtomicCommits    bool   `yaml:"atomic_commits" mapstructure:"atomic_commits"`      // default false
	TDD              bool   `yaml:"tdd" mapstructure:"tdd"`                            // default false
	TestGeneration   bool   `yaml:"test_generation" mapstructure:"test_generation"`    // default false
	ResponseLanguage string `yaml:"response_language" mapstructure:"response_language"` // e.g. "zh-TW"
	DocumentLanguage string `yaml:"document_language" mapstructure:"document_language"` // e.g. "en"
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
	}
}
