package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaults_ReturnsExpectedValues(t *testing.T) {
	d := Defaults()
	assert.Equal(t, "single", d.ExecutionMode, "default ExecutionMode should be 'single'")
	assert.Equal(t, 1, d.AgentCount, "default AgentCount should be 1")
	assert.False(t, d.AtomicCommits, "default AtomicCommits should be false")
	assert.False(t, d.TDD, "default TDD should be false")
	assert.False(t, d.TestGeneration, "default TestGeneration should be false")
	assert.Equal(t, "", d.ResponseLanguage, "default ResponseLanguage should be empty")
	assert.Equal(t, "", d.DocumentLanguage, "default DocumentLanguage should be empty")
}

func TestLoad_NoConfigFile_ReturnsDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, err := Load(tmpDir)
	require.NoError(t, err, "Load with no config file should not return error (convention over config)")
	assert.Equal(t, "single", cfg.ExecutionMode)
	assert.Equal(t, 1, cfg.AgentCount)
	assert.False(t, cfg.AtomicCommits)
	assert.False(t, cfg.TDD)
}

func TestLoad_WithConfigFile_ReturnsOverriddenValues(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	err := os.MkdirAll(claudeDir, 0755)
	require.NoError(t, err)

	configContent := `execution_mode: wave
agent_count: 4
atomic_commits: true
tdd: true
test_generation: true
response_language: zh-TW
document_language: en
`
	err = os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := Load(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, "wave", cfg.ExecutionMode)
	assert.Equal(t, 4, cfg.AgentCount)
	assert.True(t, cfg.AtomicCommits)
	assert.True(t, cfg.TDD)
	assert.True(t, cfg.TestGeneration)
	assert.Equal(t, "zh-TW", cfg.ResponseLanguage)
	assert.Equal(t, "en", cfg.DocumentLanguage)
}

func TestLoad_PartialConfigFile_MissingFieldsUseDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	err := os.MkdirAll(claudeDir, 0755)
	require.NoError(t, err)

	// Only override agent_count
	configContent := `agent_count: 3
`
	err = os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := Load(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, 3, cfg.AgentCount, "overridden value should be used")
	assert.Equal(t, "single", cfg.ExecutionMode, "non-overridden field should use default")
	assert.False(t, cfg.TDD, "non-overridden bool should use default false")
}

func TestProjectConfig_AllFields(t *testing.T) {
	// Verify all expected fields exist in the struct
	cfg := ProjectConfig{
		ExecutionMode:    "single",
		AgentCount:       1,
		AtomicCommits:    false,
		TDD:              false,
		TestGeneration:   false,
		ResponseLanguage: "zh-TW",
		DocumentLanguage: "en",
	}
	assert.NotNil(t, cfg)
}

// --- ModelProfile tests ---

func TestDefaults_ModelProfile_IsBalanced(t *testing.T) {
	d := Defaults()
	assert.Equal(t, "balanced", d.ModelProfile, "default ModelProfile should be 'balanced'")
}

func TestLoad_WithModelProfile_ReturnsQuality(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	err := os.MkdirAll(claudeDir, 0755)
	require.NoError(t, err)

	configContent := `model_profile: quality
`
	err = os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := Load(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, "quality", cfg.ModelProfile)
}

func TestResolveModel_QualityProfile_ExecutorReturnsSonnet(t *testing.T) {
	model := ResolveModel("executor", "quality", nil)
	assert.Equal(t, "claude-sonnet-4-5", model, "quality profile executor should map to sonnet")
}

func TestResolveModel_BudgetProfile_ExecutorReturnsHaiku(t *testing.T) {
	model := ResolveModel("executor", "budget", nil)
	assert.Equal(t, "claude-haiku-3-5", model, "budget profile executor should map to haiku")
}

func TestResolveModel_BalancedProfile_PlannerReturnsSonnet(t *testing.T) {
	model := ResolveModel("planner", "balanced", nil)
	assert.Equal(t, "claude-sonnet-4-5", model, "balanced profile planner should map to sonnet")
}

func TestResolveModel_WithOverride_UsesOverride(t *testing.T) {
	overrides := map[string]string{
		"executor": "claude-opus-4",
	}
	model := ResolveModel("executor", "quality", overrides)
	assert.Equal(t, "claude-opus-4", model, "model_overrides should take precedence over profile mapping")
}
