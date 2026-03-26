package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xenciscbc/mysd/internal/update"
)

// UpdateOutput is the JSON output structure for the update command.
type UpdateOutput struct {
	CurrentVersion  string      `json:"current_version"`
	LatestVersion   string      `json:"latest_version,omitempty"`
	UpdateAvailable bool        `json:"update_available"`
	ReleaseURL      string      `json:"release_url,omitempty"`
	CheckOnly       bool        `json:"check_only"`
	Force           bool        `json:"force"`
	BinaryUpdated   bool        `json:"binary_updated"`
	PluginSync      *SyncOutput `json:"plugin_sync,omitempty"`
	Error           string      `json:"error,omitempty"`
}

// SyncOutput reports plugin synchronization results.
type SyncOutput struct {
	Added   int      `json:"added"`
	Updated int      `json:"updated"`
	Deleted int      `json:"deleted"`
	Errors  []string `json:"errors,omitempty"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for updates and update mysd binary + plugins",
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Bool("check", false, "only check for updates, do not install")
	updateCmd.Flags().Bool("force", false, "skip confirmation prompt")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	checkOnly, _ := cmd.Flags().GetBool("check")
	force, _ := cmd.Flags().GetBool("force")

	currentVersion := rootCmd.Version
	if currentVersion == "" {
		currentVersion = "dev"
	}

	out := UpdateOutput{
		CurrentVersion: currentVersion,
		CheckOnly:      checkOnly,
		Force:          force,
	}

	ctx := context.Background()

	// Step 1: Check latest version from GitHub.
	// Per D-06: network failure does not block plugin sync — record error and continue.
	release, err := update.CheckLatestVersion(ctx, nil)
	if err != nil {
		out.Error = err.Error()
	} else {
		out.LatestVersion = release.TagName
		out.ReleaseURL = release.HTMLURL

		available, isAvailErr := update.IsUpdateAvailable(currentVersion, release.TagName)
		if isAvailErr != nil {
			// Non-fatal: treat as "no update" but record error
			if out.Error == "" {
				out.Error = isAvailErr.Error()
			}
		} else {
			out.UpdateAvailable = available
		}
	}

	// Step 2: If --check flag, output version info and return (D-20).
	if checkOnly {
		return printJSON(cmd, out)
	}

	// Step 3: If update is available and --force flag, apply binary update.
	if out.UpdateAvailable && force {
		exePath, exeErr := os.Executable()
		if exeErr != nil {
			out.Error = fmt.Sprintf("get executable path: %v", exeErr)
		} else {
			if applyErr := update.ApplyUpdate(ctx, nil, release, exePath); applyErr != nil {
				// Attempt rollback if update failed
				_ = update.Rollback(exePath)
				if out.Error == "" {
					out.Error = applyErr.Error()
				}
			} else {
				out.BinaryUpdated = true
			}
		}
	}

	// Step 4: Plugin sync (always runs, even if binary update failed per D-06).
	syncOut := runPluginSync(currentVersion)
	out.PluginSync = syncOut

	return printJSON(cmd, out)
}

// runPluginSync performs plugin manifest diff and sync.
// It locates the .claude/ directory relative to the executable, loads the old
// manifest, generates the new manifest from plugin/, diffs them, syncs, and
// saves the updated manifest.
func runPluginSync(version string) *SyncOutput {
	syncOut := &SyncOutput{}

	// Locate .claude/ directory — walk up from cwd to find it.
	claudeDir, err := findClaudeDir()
	if err != nil {
		syncOut.Errors = append(syncOut.Errors, fmt.Sprintf("locate .claude dir: %v", err))
		return syncOut
	}

	manifestPath := filepath.Join(claudeDir, "plugin-manifest.json")

	// Load old manifest (nil if file doesn't exist — pre-v1.1 install)
	oldManifest, err := update.LoadManifest(manifestPath)
	if err != nil {
		syncOut.Errors = append(syncOut.Errors, fmt.Sprintf("load old manifest: %v", err))
		return syncOut
	}

	// Locate plugin/ source directory (sibling of .claude/ — at project root)
	pluginDir := filepath.Join(filepath.Dir(claudeDir), "plugin")

	// Generate new manifest from the plugin/ directory
	newManifest, err := update.GenerateManifest(pluginDir, version)
	if err != nil {
		syncOut.Errors = append(syncOut.Errors, fmt.Sprintf("generate manifest: %v", err))
		return syncOut
	}

	// Diff and sync
	diff := update.DiffManifests(oldManifest, newManifest)
	result := update.SyncPlugins(pluginDir, claudeDir, diff)

	syncOut.Added = result.Added
	syncOut.Updated = result.Updated
	syncOut.Deleted = result.Deleted
	syncOut.Errors = append(syncOut.Errors, result.Errors...)

	// Save updated manifest
	if saveErr := update.SaveManifest(manifestPath, newManifest); saveErr != nil {
		syncOut.Errors = append(syncOut.Errors, fmt.Sprintf("save manifest: %v", saveErr))
	}

	return syncOut
}

// findClaudeDir searches for the .claude directory starting from the current
// working directory, walking up to the filesystem root.
func findClaudeDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %w", err)
	}

	dir := cwd
	for {
		candidate := filepath.Join(dir, ".claude")
		if info, statErr := os.Stat(candidate); statErr == nil && info.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	return "", fmt.Errorf(".claude directory not found from %s", cwd)
}

// printJSON marshals v as indented JSON and writes to cmd.OutOrStdout().
func printJSON(cmd *cobra.Command, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}
