package spec

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// OpenSpecConfig holds the project-level OpenSpec configuration (D-07).
// Locale is the source of truth for language settings (D-09).
type OpenSpecConfig struct {
	Project string `yaml:"project"`
	Locale  string `yaml:"locale"`   // BCP47: zh-TW, en-US, ja-JP (D-08)
	SpecDir string `yaml:"spec_dir"` // convention default: "openspec/specs"
	Created string `yaml:"created"`  // RFC3339 timestamp
}

// WriteOpenSpecConfig writes the OpenSpec config to {projectRoot}/openspec/config.yaml.
// Creates the openspec/ directory if it does not exist (pitfall 5).
func WriteOpenSpecConfig(projectRoot string, cfg OpenSpecConfig) error {
	dir := filepath.Join(projectRoot, "openspec")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create openspec dir: %w", err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal openspec config: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, "config.yaml"), data, 0644)
}

// ReadOpenSpecConfig reads the OpenSpec config from {projectRoot}/openspec/config.yaml.
// Returns zero-value OpenSpecConfig and nil error if file does not exist (convention-over-config).
func ReadOpenSpecConfig(projectRoot string) (OpenSpecConfig, error) {
	path := filepath.Join(projectRoot, "openspec", "config.yaml")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return OpenSpecConfig{}, nil
	}
	if err != nil {
		return OpenSpecConfig{}, fmt.Errorf("read openspec config: %w", err)
	}
	var cfg OpenSpecConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return OpenSpecConfig{}, fmt.Errorf("parse openspec config: %w", err)
	}
	return cfg, nil
}
