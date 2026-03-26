// Package scanner provides language-agnostic codebase analysis for the mysd scan command.
// It walks a project directory tree and produces structured JSON metadata
// for AI agent consumption (e.g., /mysd:scan spec generation).
package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ScanContext is the JSON-serializable output of BuildScanContext.
// It is language-agnostic — primary_language indicates the detected stack,
// and files/modules provide universal metadata for any language.
type ScanContext struct {
	RootDir         string         `json:"root_dir"`
	PrimaryLanguage string         `json:"primary_language"` // "go" | "nodejs" | "python" | "unknown"
	Files           map[string]int `json:"files"`            // extension -> count, e.g. {".go": 42}
	Modules         []ModuleInfo   `json:"modules"`          // language-agnostic module list
	ExistingSpecs   []string       `json:"existing_specs"`
	ExcludedDirs    []string       `json:"excluded_dirs"`
	TotalFiles      int            `json:"total_files"`
	ConfigExists    bool           `json:"config_exists"` // openspec/config.yaml exists?
}

// ModuleInfo contains metadata about a single module/package directory.
type ModuleInfo struct {
	Name string `json:"name"` // module/package name (relative path)
	Dir  string `json:"dir"`  // absolute path
}

// detectPrimaryLanguage checks for well-known marker files to determine the primary
// programming language of the project. Returns "go", "nodejs", "python", or "unknown".
func detectPrimaryLanguage(root string) string {
	markers := []struct {
		file string
		lang string
	}{
		{"go.mod", "go"},
		{"package.json", "nodejs"},
		{"requirements.txt", "python"},
		{"pyproject.toml", "python"},
	}
	for _, m := range markers {
		if _, err := os.Stat(filepath.Join(root, m.file)); err == nil {
			return m.lang
		}
	}
	return "unknown"
}

// BuildScanContext walks the directory tree rooted at root, collecting language-agnostic
// file and module metadata. Directories in exclude are skipped; hidden directories
// (names starting with ".") are always skipped.
//
// Returns a ScanContext with:
//   - PrimaryLanguage detected from marker files
//   - Files map of extension -> count (all files, all languages)
//   - Modules list based on the detected language
//   - ConfigExists indicating whether openspec/config.yaml is present
//   - ExistingSpecs listing change names found in openspec/changes/
func BuildScanContext(root string, exclude []string) (ScanContext, error) {
	excludeSet := make(map[string]bool, len(exclude))
	for _, d := range exclude {
		excludeSet[d] = true
	}

	primaryLanguage := detectPrimaryLanguage(root)
	files := make(map[string]int)

	// dirMap tracks which directories have source files (for module detection)
	// maps relDir -> dirInfo
	dirMap := make(map[string]*dirInfo)

	totalFiles := 0

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			name := d.Name()
			// Skip hidden dirs (e.g., .git, .specs), but NOT the root itself.
			if path != root && strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			// Skip explicitly excluded dirs (never skip root)
			if path != root && excludeSet[name] {
				return filepath.SkipDir
			}
			return nil
		}

		// Count all files by extension
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext != "" {
			files[ext]++
		}
		totalFiles++

		// Track directory membership for module detection
		absDir := filepath.Dir(path)
		relDir, err := filepath.Rel(root, absDir)
		if err != nil {
			return err
		}
		relDir = filepath.ToSlash(relDir)

		if dirMap[relDir] == nil {
			dirMap[relDir] = &dirInfo{absDir: absDir}
		}
		dirMap[relDir].fileNames = append(dirMap[relDir].fileNames, d.Name())

		return nil
	})
	if err != nil {
		return ScanContext{}, err
	}

	// Build Modules based on detected language
	modules := buildModules(primaryLanguage, dirMap)

	// Detect existing specs in openspec/changes/
	existingSpecs := detectExistingSpecs(root)

	// Check openspec/config.yaml existence
	configExists := false
	if _, err := os.Stat(filepath.Join(root, "openspec", "config.yaml")); err == nil {
		configExists = true
	}

	excludedDirs := exclude
	if excludedDirs == nil {
		excludedDirs = []string{}
	}

	return ScanContext{
		RootDir:         root,
		PrimaryLanguage: primaryLanguage,
		Files:           files,
		Modules:         modules,
		ExistingSpecs:   existingSpecs,
		ExcludedDirs:    excludedDirs,
		TotalFiles:      totalFiles,
		ConfigExists:    configExists,
	}, nil
}

// dirInfo holds per-directory metadata for module detection.
type dirInfo struct {
	absDir    string
	fileNames []string
}

// buildModules constructs the module list based on primary language and directory contents.
func buildModules(primaryLanguage string, dirMap map[string]*dirInfo) []ModuleInfo {
	modules := []ModuleInfo{}

	for relDir, info := range dirMap {
		if shouldIncludeAsModule(primaryLanguage, info.fileNames) {
			modules = append(modules, ModuleInfo{
				Name: relDir,
				Dir:  info.absDir,
			})
		}
	}

	return modules
}

// shouldIncludeAsModule returns true if a directory's files qualify it as a module
// for the given primary language.
func shouldIncludeAsModule(primaryLanguage string, fileNames []string) bool {
	fileSet := make(map[string]bool, len(fileNames))
	for _, f := range fileNames {
		fileSet[f] = true
	}

	switch primaryLanguage {
	case "go":
		// Go module: directory containing .go files
		for _, f := range fileNames {
			if strings.HasSuffix(f, ".go") {
				return true
			}
		}
	case "nodejs":
		// Node.js module: directory containing index.js/index.ts or package.json
		if fileSet["index.js"] || fileSet["index.ts"] || fileSet["package.json"] {
			return true
		}
	case "python":
		// Python module: directory containing __init__.py
		if fileSet["__init__.py"] {
			return true
		}
	default:
		// Unknown: include top-level directories that have source files
		return len(fileNames) > 0
	}

	return false
}

// detectExistingSpecs finds change names in openspec/changes/ directory.
func detectExistingSpecs(root string) []string {
	specs := []string{}
	changesDir := filepath.Join(root, "openspec", "changes")

	entries, err := os.ReadDir(changesDir)
	if err != nil {
		// openspec/changes/ not present — normal for first-time scan
		return specs
	}

	for _, e := range entries {
		if e.IsDir() {
			specs = append(specs, e.Name())
		}
	}

	return specs
}
