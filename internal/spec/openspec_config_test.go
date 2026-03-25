package spec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWriteOpenSpecConfig_CreatesFile verifies WriteOpenSpecConfig creates openspec/config.yaml.
func TestWriteOpenSpecConfig_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	cfg := OpenSpecConfig{
		Project: "my-project",
		Locale:  "en-US",
		SpecDir: "openspec/specs",
		Created: "2026-03-25T00:00:00Z",
	}

	err := WriteOpenSpecConfig(dir, cfg)
	require.NoError(t, err)

	configPath := filepath.Join(dir, "openspec", "config.yaml")
	assert.FileExists(t, configPath)
}

// TestWriteOpenSpecConfig_CreatesDir verifies openspec/ directory is created if it does not exist.
func TestWriteOpenSpecConfig_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	// Write to a nested dir that does not exist yet
	cfg := OpenSpecConfig{Project: "test"}

	err := WriteOpenSpecConfig(dir, cfg)
	require.NoError(t, err)

	openspecDir := filepath.Join(dir, "openspec")
	info, err := os.Stat(openspecDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir(), "openspec/ should be a directory")
}

// TestReadOpenSpecConfig_RoundTrip verifies write then read returns all fields intact.
func TestReadOpenSpecConfig_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	want := OpenSpecConfig{
		Project: "mysd",
		Locale:  "zh-TW",
		SpecDir: "openspec/specs",
		Created: "2026-03-25T10:00:00Z",
	}

	require.NoError(t, WriteOpenSpecConfig(dir, want))

	got, err := ReadOpenSpecConfig(dir)
	require.NoError(t, err)
	assert.Equal(t, want.Project, got.Project)
	assert.Equal(t, want.Locale, got.Locale)
	assert.Equal(t, want.SpecDir, got.SpecDir)
	assert.Equal(t, want.Created, got.Created)
}

// TestReadOpenSpecConfig_NotExist verifies reading from empty dir returns zero-value and nil error.
func TestReadOpenSpecConfig_NotExist(t *testing.T) {
	dir := t.TempDir()

	cfg, err := ReadOpenSpecConfig(dir)
	require.NoError(t, err)
	assert.Equal(t, OpenSpecConfig{}, cfg, "should return zero-value when file does not exist")
}

// TestReadOpenSpecConfig_MalformedYAML verifies malformed YAML returns error.
func TestReadOpenSpecConfig_MalformedYAML(t *testing.T) {
	dir := t.TempDir()
	openspecDir := filepath.Join(dir, "openspec")
	require.NoError(t, os.MkdirAll(openspecDir, 0755))

	// Write invalid YAML
	err := os.WriteFile(filepath.Join(openspecDir, "config.yaml"), []byte(":\tbad yaml\t:\n"), 0644)
	require.NoError(t, err)

	_, err = ReadOpenSpecConfig(dir)
	assert.Error(t, err, "malformed YAML should return error")
}

// TestWriteOpenSpecConfig_BCP47Locale verifies locale field is written correctly.
func TestWriteOpenSpecConfig_BCP47Locale(t *testing.T) {
	dir := t.TempDir()
	cfg := OpenSpecConfig{
		Project: "test",
		Locale:  "zh-TW",
	}

	require.NoError(t, WriteOpenSpecConfig(dir, cfg))

	raw, err := os.ReadFile(filepath.Join(dir, "openspec", "config.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(raw), "locale: zh-TW")
}
