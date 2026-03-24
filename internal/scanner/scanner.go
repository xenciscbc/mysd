// Package scanner provides codebase analysis for the mysd scan command.
// It walks a Go project directory tree and produces structured JSON metadata
// for AI agent consumption (e.g., /mysd:scan spec generation).
package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mysd/internal/spec"
)

// ScanContext is the JSON-serializable output of BuildScanContext.
// It contains all Go package metadata for the scanned codebase.
type ScanContext struct {
	RootDir       string        `json:"root_dir"`
	Packages      []PackageInfo `json:"packages"`
	ExistingSpecs []string      `json:"existing_specs"`
	ExcludedDirs  []string      `json:"excluded_dirs"`
	TotalFiles    int           `json:"total_files"`
}

// PackageInfo contains metadata about a single Go package directory.
type PackageInfo struct {
	Name      string   `json:"name"`
	Dir       string   `json:"dir"`
	GoFiles   []string `json:"go_files"`
	TestFiles []string `json:"test_files"`
	HasSpec   bool     `json:"has_spec"`
}

// BuildScanContext walks the directory tree rooted at root, collecting Go package
// information. Directories in exclude are skipped; hidden directories (names starting
// with ".") are always skipped.
//
// Returns a ScanContext with all packages found, with HasSpec=true for packages
// that have an existing spec under the detected specs directory.
func BuildScanContext(root string, exclude []string) (ScanContext, error) {
	excludeSet := make(map[string]bool, len(exclude))
	for _, d := range exclude {
		excludeSet[d] = true
	}

	// pkgFiles maps relative package dir (forward slashes) -> {goFiles, testFiles}
	type pkgFiles struct {
		goFiles   []string
		testFiles []string
		absDir    string
	}
	pkgMap := make(map[string]*pkgFiles)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			name := d.Name()
			// Skip hidden dirs (e.g., .git, .specs), but NOT the root itself.
			// WalkDir calls the root with name "." which starts with "." —
			// we must not skip it or the entire walk is aborted.
			if path != root && strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			// Skip explicitly excluded dirs (never skip root)
			if path != root && excludeSet[name] {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .go files
		if !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}

		absDir := filepath.Dir(path)
		relDir, err := filepath.Rel(root, absDir)
		if err != nil {
			return err
		}
		// Normalize to forward slashes for cross-platform consistency
		relDir = filepath.ToSlash(relDir)
		// Use "." for root-level files; but root-level Go packages are rare
		// in multi-package projects. Keep as-is.

		if pkgMap[relDir] == nil {
			pkgMap[relDir] = &pkgFiles{absDir: absDir}
		}

		fileName := d.Name()
		if strings.HasSuffix(fileName, "_test.go") {
			pkgMap[relDir].testFiles = append(pkgMap[relDir].testFiles, fileName)
		} else {
			pkgMap[relDir].goFiles = append(pkgMap[relDir].goFiles, fileName)
		}
		return nil
	})
	if err != nil {
		return ScanContext{}, err
	}

	// Build PackageInfo list (unsorted — order follows WalkDir which is lexical)
	packages := make([]PackageInfo, 0, len(pkgMap))
	for relDir, pf := range pkgMap {
		goFiles := pf.goFiles
		if goFiles == nil {
			goFiles = []string{}
		}
		testFiles := pf.testFiles
		if testFiles == nil {
			testFiles = []string{}
		}
		packages = append(packages, PackageInfo{
			Name:      relDir,
			Dir:       pf.absDir,
			GoFiles:   goFiles,
			TestFiles: testFiles,
		})
	}

	// Detect existing specs directory for HasSpec detection
	var specsDir string
	specsDir, _, specsErr := spec.DetectSpecDir(root)
	if specsErr != nil {
		// No specs dir — normal for first-time scan; all HasSpec=false
		specsDir = ""
	}

	var existingSpecs []string

	if specsDir != "" {
		absSpecsDir := filepath.Join(root, specsDir)
		changesDir := filepath.Join(absSpecsDir, "changes")

		for i := range packages {
			pkgName := packages[i].Name
			specPath := filepath.Join(changesDir, pkgName)
			if info, err := os.Stat(specPath); err == nil && info.IsDir() {
				packages[i].HasSpec = true
				existingSpecs = append(existingSpecs, pkgName)
			}
		}
	}

	// Count total files
	totalFiles := 0
	for _, pkg := range packages {
		totalFiles += len(pkg.GoFiles) + len(pkg.TestFiles)
	}

	excludedDirs := exclude
	if excludedDirs == nil {
		excludedDirs = []string{}
	}

	return ScanContext{
		RootDir:       root,
		Packages:      packages,
		ExistingSpecs: existingSpecs,
		ExcludedDirs:  excludedDirs,
		TotalFiles:    totalFiles,
	}, nil
}
