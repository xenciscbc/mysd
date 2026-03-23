package spec

import (
	"os"
	"path/filepath"
)

// DetectSpecDir auto-detects which spec directory convention the project uses.
// It checks for .specs/ first (FlavorMySD), then openspec/ (FlavorOpenSpec).
// Returns ErrNoSpecDir if neither is found.
func DetectSpecDir(root string) (dir string, flavor SpecDirFlavor, err error) {
	mysdPath := filepath.Join(root, ".specs")
	if info, statErr := os.Stat(mysdPath); statErr == nil && info.IsDir() {
		return ".specs", FlavorMySD, nil
	}

	openspecPath := filepath.Join(root, "openspec")
	if info, statErr := os.Stat(openspecPath); statErr == nil && info.IsDir() {
		return "openspec", FlavorOpenSpec, nil
	}

	return "", FlavorNone, ErrNoSpecDir
}

// ListChanges returns the names of all change directories under specDir/changes/.
func ListChanges(specDir string) ([]string, error) {
	changesDir := filepath.Join(specDir, "changes")
	entries, err := os.ReadDir(changesDir)
	if err != nil {
		return nil, err
	}

	var changes []string
	for _, entry := range entries {
		if entry.IsDir() {
			changes = append(changes, entry.Name())
		}
	}
	return changes, nil
}
