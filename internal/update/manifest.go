// Package update provides plugin manifest tracking and synchronization for mysd self-update.
package update

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PluginManifest tracks official plugin files shipped with a release.
type PluginManifest struct {
	Version  string   `json:"version"`
	Commands []string `json:"commands"` // filenames like "mysd-update.md"
	Agents   []string `json:"agents"`   // filenames like "mysd-executor.md"
}

// ManifestDiff describes what sync operations are needed.
type ManifestDiff struct {
	AddCommands    []string // files to add to .claude/commands/
	UpdateCommands []string // files to overwrite in .claude/commands/
	DeleteCommands []string // files to remove from .claude/commands/
	AddAgents      []string
	UpdateAgents   []string
	DeleteAgents   []string
}

// LoadManifest reads a JSON plugin manifest from path.
// Returns (nil, nil) if the file does not exist — representing a pre-v1.1 installation
// without a manifest (convention-over-config, same pattern as deferred.go).
func LoadManifest(path string) (*PluginManifest, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("manifest: %w", err)
	}
	var m PluginManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("manifest: %w", err)
	}
	return &m, nil
}

// SaveManifest writes the manifest to path as indented JSON.
func SaveManifest(path string, m *PluginManifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("manifest: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("manifest: %w", err)
	}
	return nil
}

// DiffManifests computes the diff between old and new manifests.
//
// If old is nil (pre-v1.1 installation without manifest, per D-17):
//   - All files in new are "add" operations
//   - Zero "delete" operations (backward-compat: never delete when no manifest exists)
//
// If both are non-nil:
//   - Files in new but not in old → "add"
//   - Files in both → "update"
//   - Files in old but not in new → "delete"
func DiffManifests(old, newM *PluginManifest) ManifestDiff {
	var diff ManifestDiff

	if old == nil {
		// Pre-v1.1 backward compatibility: add all new files, never delete
		diff.AddCommands = append(diff.AddCommands, newM.Commands...)
		diff.AddAgents = append(diff.AddAgents, newM.Agents...)
		return diff
	}

	// Compute command diff
	diff.AddCommands, diff.UpdateCommands, diff.DeleteCommands = diffSlices(old.Commands, newM.Commands)
	// Compute agent diff
	diff.AddAgents, diff.UpdateAgents, diff.DeleteAgents = diffSlices(old.Agents, newM.Agents)

	return diff
}

// diffSlices computes add/update/delete operations between old and new string slices.
func diffSlices(old, newItems []string) (add, update, delete []string) {
	oldSet := make(map[string]struct{}, len(old))
	for _, f := range old {
		oldSet[f] = struct{}{}
	}

	newSet := make(map[string]struct{}, len(newItems))
	for _, f := range newItems {
		newSet[f] = struct{}{}
	}

	for _, f := range newItems {
		if _, exists := oldSet[f]; exists {
			update = append(update, f)
		} else {
			add = append(add, f)
		}
	}

	for _, f := range old {
		if _, exists := newSet[f]; !exists {
			delete = append(delete, f)
		}
	}

	return add, update, delete
}

// GenerateManifest scans pluginDir/commands/*.md and pluginDir/agents/*.md,
// returning a PluginManifest with all .md filenames (excluding CLAUDE.md).
func GenerateManifest(pluginDir, version string) (*PluginManifest, error) {
	commands, err := scanMDFiles(filepath.Join(pluginDir, "commands"))
	if err != nil {
		return nil, fmt.Errorf("manifest generate: %w", err)
	}
	agents, err := scanMDFiles(filepath.Join(pluginDir, "agents"))
	if err != nil {
		return nil, fmt.Errorf("manifest generate: %w", err)
	}

	return &PluginManifest{
		Version:  version,
		Commands: commands,
		Agents:   agents,
	}, nil
}

// scanMDFiles returns filenames (not full paths) of all .md files in dir,
// excluding CLAUDE.md files.
func scanMDFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return []string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".md") {
			continue
		}
		if name == "CLAUDE.md" {
			continue
		}
		files = append(files, name)
	}
	return files, nil
}
