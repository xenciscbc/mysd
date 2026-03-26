package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadManifest verifies LoadManifest behavior for missing and valid files.
func TestLoadManifest(t *testing.T) {
	t.Run("missing file returns nil manifest and nil error", func(t *testing.T) {
		m, err := LoadManifest("/nonexistent/path/plugin-manifest.json")
		assert.NoError(t, err)
		assert.Nil(t, m)
	})

	t.Run("valid JSON file returns populated manifest", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "plugin-manifest.json")
		manifest := &PluginManifest{
			Version:  "v1.1.0",
			Commands: []string{"mysd-update.md", "mysd-note.md"},
			Agents:   []string{"mysd-executor.md"},
		}
		data, err := json.MarshalIndent(manifest, "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(path, data, 0644))

		loaded, err := LoadManifest(path)
		require.NoError(t, err)
		require.NotNil(t, loaded)
		assert.Equal(t, "v1.1.0", loaded.Version)
		assert.Equal(t, []string{"mysd-update.md", "mysd-note.md"}, loaded.Commands)
		assert.Equal(t, []string{"mysd-executor.md"}, loaded.Agents)
	})
}

// TestSaveManifest verifies SaveManifest writes indented JSON to disk.
func TestSaveManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "plugin-manifest.json")
	manifest := &PluginManifest{
		Version:  "v1.2.0",
		Commands: []string{"mysd-apply.md"},
		Agents:   []string{"mysd-planner.md"},
	}

	err := SaveManifest(path, manifest)
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	// Verify it's valid JSON and round-trips correctly
	var loaded PluginManifest
	require.NoError(t, json.Unmarshal(data, &loaded))
	assert.Equal(t, manifest.Version, loaded.Version)
	assert.Equal(t, manifest.Commands, loaded.Commands)
	assert.Equal(t, manifest.Agents, loaded.Agents)

	// Verify it's indented (contains newlines)
	assert.Contains(t, string(data), "\n")
}

// TestDiffManifests verifies DiffManifests computes correct add/update/delete operations.
func TestDiffManifests(t *testing.T) {
	t.Run("nil old manifest: all new files are add, zero deletes (backward compat D-17)", func(t *testing.T) {
		newManifest := &PluginManifest{
			Version:  "v1.1.0",
			Commands: []string{"mysd-update.md", "mysd-note.md"},
			Agents:   []string{"mysd-executor.md", "mysd-planner.md"},
		}
		diff := DiffManifests(nil, newManifest)

		assert.ElementsMatch(t, []string{"mysd-update.md", "mysd-note.md"}, diff.AddCommands)
		assert.Empty(t, diff.UpdateCommands)
		assert.Empty(t, diff.DeleteCommands)
		assert.ElementsMatch(t, []string{"mysd-executor.md", "mysd-planner.md"}, diff.AddAgents)
		assert.Empty(t, diff.UpdateAgents)
		assert.Empty(t, diff.DeleteAgents)
	})

	t.Run("new has extra file: that file is add", func(t *testing.T) {
		old := &PluginManifest{
			Version:  "v1.0.0",
			Commands: []string{"mysd-note.md"},
			Agents:   []string{},
		}
		newM := &PluginManifest{
			Version:  "v1.1.0",
			Commands: []string{"mysd-note.md", "mysd-update.md"},
			Agents:   []string{},
		}
		diff := DiffManifests(old, newM)

		assert.Contains(t, diff.AddCommands, "mysd-update.md")
		assert.NotContains(t, diff.AddCommands, "mysd-note.md")
		assert.Empty(t, diff.DeleteCommands)
	})

	t.Run("old has extra file: that file is delete", func(t *testing.T) {
		old := &PluginManifest{
			Version:  "v1.0.0",
			Commands: []string{"mysd-note.md", "mysd-old.md"},
			Agents:   []string{},
		}
		newM := &PluginManifest{
			Version:  "v1.1.0",
			Commands: []string{"mysd-note.md"},
			Agents:   []string{},
		}
		diff := DiffManifests(old, newM)

		assert.Contains(t, diff.DeleteCommands, "mysd-old.md")
		assert.NotContains(t, diff.DeleteCommands, "mysd-note.md")
		assert.Empty(t, diff.AddCommands)
	})

	t.Run("both have same file: that file is update", func(t *testing.T) {
		old := &PluginManifest{
			Version:  "v1.0.0",
			Commands: []string{"mysd-note.md"},
			Agents:   []string{"mysd-executor.md"},
		}
		newM := &PluginManifest{
			Version:  "v1.1.0",
			Commands: []string{"mysd-note.md"},
			Agents:   []string{"mysd-executor.md"},
		}
		diff := DiffManifests(old, newM)

		assert.Contains(t, diff.UpdateCommands, "mysd-note.md")
		assert.Empty(t, diff.AddCommands)
		assert.Empty(t, diff.DeleteCommands)
		assert.Contains(t, diff.UpdateAgents, "mysd-executor.md")
		assert.Empty(t, diff.AddAgents)
		assert.Empty(t, diff.DeleteAgents)
	})

	t.Run("identical files: empty diff", func(t *testing.T) {
		old := &PluginManifest{
			Version:  "v1.0.0",
			Commands: []string{"mysd-note.md"},
			Agents:   []string{},
		}
		newM := &PluginManifest{
			Version:  "v1.0.0",
			Commands: []string{"mysd-note.md"},
			Agents:   []string{},
		}
		diff := DiffManifests(old, newM)

		assert.Empty(t, diff.AddCommands)
		assert.Empty(t, diff.DeleteCommands)
		assert.Contains(t, diff.UpdateCommands, "mysd-note.md")
	})

	t.Run("agents: add, update, delete work independently", func(t *testing.T) {
		old := &PluginManifest{
			Version:  "v1.0.0",
			Commands: []string{},
			Agents:   []string{"mysd-executor.md", "mysd-old-agent.md"},
		}
		newM := &PluginManifest{
			Version:  "v1.1.0",
			Commands: []string{},
			Agents:   []string{"mysd-executor.md", "mysd-new-agent.md"},
		}
		diff := DiffManifests(old, newM)

		assert.Contains(t, diff.AddAgents, "mysd-new-agent.md")
		assert.Contains(t, diff.UpdateAgents, "mysd-executor.md")
		assert.Contains(t, diff.DeleteAgents, "mysd-old-agent.md")
	})
}

// TestGenerateManifest verifies GenerateManifest scans directories and returns manifest.
func TestGenerateManifest(t *testing.T) {
	dir := t.TempDir()

	// Create plugin structure
	commandsDir := filepath.Join(dir, "commands")
	agentsDir := filepath.Join(dir, "agents")
	require.NoError(t, os.MkdirAll(commandsDir, 0755))
	require.NoError(t, os.MkdirAll(agentsDir, 0755))

	// Create .md files
	require.NoError(t, os.WriteFile(filepath.Join(commandsDir, "mysd-update.md"), []byte("# update"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(commandsDir, "mysd-note.md"), []byte("# note"), 0644))
	// CLAUDE.md should be excluded
	require.NoError(t, os.WriteFile(filepath.Join(commandsDir, "CLAUDE.md"), []byte("# claude"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(agentsDir, "mysd-executor.md"), []byte("# executor"), 0644))

	manifest, err := GenerateManifest(dir, "v1.1.0")
	require.NoError(t, err)
	require.NotNil(t, manifest)

	assert.Equal(t, "v1.1.0", manifest.Version)
	assert.ElementsMatch(t, []string{"mysd-update.md", "mysd-note.md"}, manifest.Commands)
	assert.ElementsMatch(t, []string{"mysd-executor.md"}, manifest.Agents)
	// CLAUDE.md must not be included
	assert.NotContains(t, manifest.Commands, "CLAUDE.md")
}
