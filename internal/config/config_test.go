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
	assert.Equal(t, "sonnet", model, "quality profile executor should map to sonnet")
}

func TestResolveModel_BudgetProfile_ExecutorReturnsHaiku(t *testing.T) {
	model := ResolveModel("executor", "budget", nil)
	assert.Equal(t, "haiku", model, "budget profile executor should map to haiku")
}

func TestResolveModel_BalancedProfile_PlannerReturnsOpus(t *testing.T) {
	model := ResolveModel("planner", "balanced", nil)
	assert.Equal(t, "opus", model, "balanced profile planner should map to opus")
}

func TestResolveModel_WithOverride_UsesOverride(t *testing.T) {
	overrides := map[string]string{
		"executor": "opus",
	}
	model := ResolveModel("executor", "quality", overrides)
	assert.Equal(t, "opus", model, "model_overrides should take precedence over profile mapping")
}

// TestResolveModel_AllRoles verifies all 11 agent roles return correct model across all 3 profiles.
func TestResolveModel_AllRoles(t *testing.T) {
	expected := map[string]map[string]string{
		"quality": {
			"spec-writer": "opus", "designer": "opus", "planner": "opus",
			"executor": "sonnet", "spec-executor": "opus", "verifier": "opus", "fast-forward": "sonnet",
			"researcher": "opus", "advisor": "opus", "proposal-writer": "opus", "plan-checker": "opus",
			"reviewer": "opus",
		},
		"balanced": {
			"spec-writer": "opus", "designer": "opus", "planner": "opus",
			"executor": "sonnet", "spec-executor": "opus", "verifier": "opus", "fast-forward": "sonnet",
			"researcher": "sonnet", "advisor": "opus", "proposal-writer": "sonnet", "plan-checker": "opus",
			"reviewer": "sonnet",
		},
		"budget": {
			"spec-writer": "sonnet", "designer": "haiku", "planner": "sonnet",
			"executor": "haiku", "spec-executor": "sonnet", "verifier": "sonnet", "fast-forward": "haiku",
			"researcher": "sonnet", "advisor": "sonnet", "proposal-writer": "sonnet", "plan-checker": "sonnet",
			"reviewer": "sonnet",
		},
	}

	for profile, roles := range expected {
		for role, want := range roles {
			t.Run(role+"/"+profile, func(t *testing.T) {
				got := ResolveModel(role, profile, nil)
				assert.Equal(t, want, got,
					"role %s profile %s should map to %s", role, profile, want)
			})
		}
	}
}

// TestResolveModel_SpecExecutorRole verifies spec-executor role returns correct model per profile.
func TestResolveModel_SpecExecutorRole(t *testing.T) {
	assert.Equal(t, "opus", ResolveModel("spec-executor", "quality", nil), "quality spec-executor should use opus")
	assert.Equal(t, "opus", ResolveModel("spec-executor", "balanced", nil), "balanced spec-executor should use opus")
	assert.Equal(t, "sonnet", ResolveModel("spec-executor", "budget", nil), "budget spec-executor should use sonnet")
}

// TestResolveModel_ReviewerRole verifies reviewer role returns correct model per profile.
func TestResolveModel_ReviewerRole(t *testing.T) {
	assert.Equal(t, "opus", ResolveModel("reviewer", "quality", nil), "quality reviewer should use opus")
	assert.Equal(t, "sonnet", ResolveModel("reviewer", "balanced", nil), "balanced reviewer should use sonnet")
	assert.Equal(t, "sonnet", ResolveModel("reviewer", "budget", nil), "budget reviewer should use sonnet")
}

// TestResolveModel_NewRoles_Override verifies overrides work for new roles.
func TestResolveModel_NewRoles_Override(t *testing.T) {
	overrides := map[string]string{
		"plan-checker": "custom-model",
	}
	model := ResolveModel("plan-checker", "quality", overrides)
	assert.Equal(t, "custom-model", model, "override should take precedence for new roles")
}

// TestLoad_WithSpecExecutionMode verifies execution_mode: "spec" is accepted in config parsing.
func TestLoad_WithSpecExecutionMode(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	err := os.MkdirAll(claudeDir, 0755)
	require.NoError(t, err)

	configContent := `execution_mode: spec
`
	err = os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := Load(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, "spec", cfg.ExecutionMode)
}

// TestDefaults_NewFields verifies new ProjectConfig fields have correct defaults.
func TestDefaults_NewFields(t *testing.T) {
	d := Defaults()
	assert.Equal(t, ".worktrees", d.WorktreeDir, "default WorktreeDir should be '.worktrees'")
	assert.False(t, d.AutoMode, "default AutoMode should be false")
}

// TestLoad_NewFields verifies worktree_dir and auto_mode are read from config file.
func TestLoad_NewFields(t *testing.T) {
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	err := os.MkdirAll(claudeDir, 0755)
	require.NoError(t, err)

	configContent := `worktree_dir: .wt
auto_mode: true
`
	err = os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := Load(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, ".wt", cfg.WorktreeDir)
	assert.True(t, cfg.AutoMode)
}
