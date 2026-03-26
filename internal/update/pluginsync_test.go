package update

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSyncPlugins verifies SyncPlugins copies and deletes files based on ManifestDiff.
func TestSyncPlugins(t *testing.T) {
	t.Run("empty diff returns all-zero SyncResult", func(t *testing.T) {
		sourceDir := t.TempDir()
		targetDir := t.TempDir()

		diff := ManifestDiff{}
		result := SyncPlugins(sourceDir, targetDir, diff)

		assert.Equal(t, 0, result.Added)
		assert.Equal(t, 0, result.Updated)
		assert.Equal(t, 0, result.Deleted)
		assert.Empty(t, result.Errors)
	})

	t.Run("add command: copies file from source to target commands dir", func(t *testing.T) {
		sourceDir := t.TempDir()
		targetDir := t.TempDir()

		// Create source file
		require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "commands"), 0755))
		require.NoError(t, os.WriteFile(
			filepath.Join(sourceDir, "commands", "mysd-update.md"),
			[]byte("# update command"),
			0644,
		))

		diff := ManifestDiff{
			AddCommands: []string{"mysd-update.md"},
		}
		result := SyncPlugins(sourceDir, targetDir, diff)

		assert.Equal(t, 1, result.Added)
		assert.Equal(t, 0, result.Updated)
		assert.Equal(t, 0, result.Deleted)
		assert.Empty(t, result.Errors)

		// Verify file was copied
		content, err := os.ReadFile(filepath.Join(targetDir, "commands", "mysd-update.md"))
		require.NoError(t, err)
		assert.Equal(t, "# update command", string(content))
	})

	t.Run("update command: overwrites existing file", func(t *testing.T) {
		sourceDir := t.TempDir()
		targetDir := t.TempDir()

		// Create existing target file
		require.NoError(t, os.MkdirAll(filepath.Join(targetDir, "commands"), 0755))
		require.NoError(t, os.WriteFile(
			filepath.Join(targetDir, "commands", "mysd-note.md"),
			[]byte("# old content"),
			0644,
		))

		// Create new source file
		require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "commands"), 0755))
		require.NoError(t, os.WriteFile(
			filepath.Join(sourceDir, "commands", "mysd-note.md"),
			[]byte("# new content"),
			0644,
		))

		diff := ManifestDiff{
			UpdateCommands: []string{"mysd-note.md"},
		}
		result := SyncPlugins(sourceDir, targetDir, diff)

		assert.Equal(t, 0, result.Added)
		assert.Equal(t, 1, result.Updated)
		assert.Equal(t, 0, result.Deleted)
		assert.Empty(t, result.Errors)

		// Verify file was overwritten
		content, err := os.ReadFile(filepath.Join(targetDir, "commands", "mysd-note.md"))
		require.NoError(t, err)
		assert.Equal(t, "# new content", string(content))
	})

	t.Run("delete command: removes file from target", func(t *testing.T) {
		sourceDir := t.TempDir()
		targetDir := t.TempDir()

		// Create existing target file to be deleted
		require.NoError(t, os.MkdirAll(filepath.Join(targetDir, "commands"), 0755))
		require.NoError(t, os.WriteFile(
			filepath.Join(targetDir, "commands", "mysd-old.md"),
			[]byte("# old command"),
			0644,
		))

		diff := ManifestDiff{
			DeleteCommands: []string{"mysd-old.md"},
		}
		result := SyncPlugins(sourceDir, targetDir, diff)

		assert.Equal(t, 0, result.Added)
		assert.Equal(t, 0, result.Updated)
		assert.Equal(t, 1, result.Deleted)
		assert.Empty(t, result.Errors)

		// Verify file was deleted
		_, err := os.Stat(filepath.Join(targetDir, "commands", "mysd-old.md"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("creates target directories if they do not exist", func(t *testing.T) {
		sourceDir := t.TempDir()
		targetDir := t.TempDir()

		// Source has the file but target commands/ dir doesn't exist
		require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "commands"), 0755))
		require.NoError(t, os.WriteFile(
			filepath.Join(sourceDir, "commands", "mysd-new.md"),
			[]byte("# new"),
			0644,
		))

		diff := ManifestDiff{
			AddCommands: []string{"mysd-new.md"},
		}
		result := SyncPlugins(sourceDir, targetDir, diff)

		assert.Equal(t, 1, result.Added)
		assert.Empty(t, result.Errors)

		// Target commands dir should have been created
		_, err := os.Stat(filepath.Join(targetDir, "commands"))
		assert.NoError(t, err)
	})

	t.Run("add agent: copies file from source to target agents dir", func(t *testing.T) {
		sourceDir := t.TempDir()
		targetDir := t.TempDir()

		require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "agents"), 0755))
		require.NoError(t, os.WriteFile(
			filepath.Join(sourceDir, "agents", "mysd-executor.md"),
			[]byte("# executor agent"),
			0644,
		))

		diff := ManifestDiff{
			AddAgents: []string{"mysd-executor.md"},
		}
		result := SyncPlugins(sourceDir, targetDir, diff)

		assert.Equal(t, 1, result.Added)
		assert.Equal(t, 0, result.Deleted)
		assert.Empty(t, result.Errors)

		content, err := os.ReadFile(filepath.Join(targetDir, "agents", "mysd-executor.md"))
		require.NoError(t, err)
		assert.Equal(t, "# executor agent", string(content))
	})

	t.Run("delete agent: removes file from target agents dir", func(t *testing.T) {
		sourceDir := t.TempDir()
		targetDir := t.TempDir()

		require.NoError(t, os.MkdirAll(filepath.Join(targetDir, "agents"), 0755))
		require.NoError(t, os.WriteFile(
			filepath.Join(targetDir, "agents", "mysd-old-agent.md"),
			[]byte("# old agent"),
			0644,
		))

		diff := ManifestDiff{
			DeleteAgents: []string{"mysd-old-agent.md"},
		}
		result := SyncPlugins(sourceDir, targetDir, diff)

		assert.Equal(t, 1, result.Deleted)
		assert.Empty(t, result.Errors)

		_, err := os.Stat(filepath.Join(targetDir, "agents", "mysd-old-agent.md"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("delete non-existent file: non-fatal error appended", func(t *testing.T) {
		sourceDir := t.TempDir()
		targetDir := t.TempDir()

		diff := ManifestDiff{
			DeleteCommands: []string{"does-not-exist.md"},
		}
		result := SyncPlugins(sourceDir, targetDir, diff)

		// Delete errors are non-fatal — appended to Errors slice, not returned as error
		assert.Equal(t, 0, result.Deleted)
		assert.NotEmpty(t, result.Errors)
	})

	t.Run("combined add update delete: counts are correct", func(t *testing.T) {
		sourceDir := t.TempDir()
		targetDir := t.TempDir()

		require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "commands"), 0755))
		require.NoError(t, os.MkdirAll(filepath.Join(targetDir, "commands"), 0755))

		// File to add (only in source)
		require.NoError(t, os.WriteFile(
			filepath.Join(sourceDir, "commands", "mysd-new.md"),
			[]byte("# new"),
			0644,
		))
		// File to update (in both)
		require.NoError(t, os.WriteFile(
			filepath.Join(sourceDir, "commands", "mysd-update.md"),
			[]byte("# updated"),
			0644,
		))
		require.NoError(t, os.WriteFile(
			filepath.Join(targetDir, "commands", "mysd-update.md"),
			[]byte("# old"),
			0644,
		))
		// File to delete (only in target)
		require.NoError(t, os.WriteFile(
			filepath.Join(targetDir, "commands", "mysd-old.md"),
			[]byte("# old"),
			0644,
		))

		diff := ManifestDiff{
			AddCommands:    []string{"mysd-new.md"},
			UpdateCommands: []string{"mysd-update.md"},
			DeleteCommands: []string{"mysd-old.md"},
		}
		result := SyncPlugins(sourceDir, targetDir, diff)

		assert.Equal(t, 1, result.Added)
		assert.Equal(t, 1, result.Updated)
		assert.Equal(t, 1, result.Deleted)
		assert.Empty(t, result.Errors)
	})
}
