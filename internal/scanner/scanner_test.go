package scanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xenciscbc/mysd/internal/scanner"
)

// writeFile creates a file with given content in a temp directory.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
}

func TestBuildScanContext_GoProject(t *testing.T) {
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "go.mod"), "module example.com/myapp\n\ngo 1.21\n")
	writeFile(t, filepath.Join(root, "cmd", "main.go"), "package main\n")
	writeFile(t, filepath.Join(root, "internal", "foo", "foo.go"), "package foo\n")
	writeFile(t, filepath.Join(root, "internal", "foo", "foo_test.go"), "package foo_test\n")

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	assert.Equal(t, root, ctx.RootDir)
	assert.Equal(t, "go", ctx.PrimaryLanguage)
	assert.Contains(t, ctx.Files, ".go", "Files map should have .go key")
	assert.GreaterOrEqual(t, ctx.Files[".go"], 3, "should have at least 3 .go files")
	assert.NotNil(t, ctx.Modules, "Modules must not be nil")
	assert.GreaterOrEqual(t, len(ctx.Modules), 1, "Go project should detect at least one module")
}

func TestBuildScanContext_NodeProject(t *testing.T) {
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "package.json"), `{"name":"myapp","version":"1.0.0"}`)
	writeFile(t, filepath.Join(root, "src", "index.js"), "console.log('hello');\n")
	writeFile(t, filepath.Join(root, "src", "app.ts"), "export const x = 1;\n")

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	assert.Equal(t, "nodejs", ctx.PrimaryLanguage)
	assert.Contains(t, ctx.Files, ".js", "Files map should have .js key")
	assert.Contains(t, ctx.Files, ".ts", "Files map should have .ts key")
}

func TestBuildScanContext_PythonProject(t *testing.T) {
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "pyproject.toml"), "[project]\nname = \"myapp\"\n")
	writeFile(t, filepath.Join(root, "myapp", "__init__.py"), "")
	writeFile(t, filepath.Join(root, "myapp", "main.py"), "print('hello')\n")

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	assert.Equal(t, "python", ctx.PrimaryLanguage)
	assert.Contains(t, ctx.Files, ".py", "Files map should have .py key")
}

func TestBuildScanContext_UnknownProject(t *testing.T) {
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "readme.txt"), "just a text file\n")
	writeFile(t, filepath.Join(root, "notes.md"), "# Notes\n")

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	assert.Equal(t, "unknown", ctx.PrimaryLanguage)
}

func TestBuildScanContext_ConfigExists(t *testing.T) {
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "go.mod"), "module example.com/myapp\n\ngo 1.21\n")
	writeFile(t, filepath.Join(root, "openspec", "config.yaml"), "project: myapp\nlocale: en-US\n")

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	assert.True(t, ctx.ConfigExists, "config_exists should be true when openspec/config.yaml exists")
}

func TestBuildScanContext_ConfigNotExists(t *testing.T) {
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "go.mod"), "module example.com/myapp\n\ngo 1.21\n")

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	assert.False(t, ctx.ConfigExists, "config_exists should be false when openspec/config.yaml absent")
}

func TestBuildScanContext_ExcludeDirs(t *testing.T) {
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "go.mod"), "module example.com/myapp\n\ngo 1.21\n")
	writeFile(t, filepath.Join(root, "vendor", "lib", "lib.go"), "package lib\n")
	writeFile(t, filepath.Join(root, "internal", "bar", "bar.go"), "package bar\n")

	ctx, err := scanner.BuildScanContext(root, []string{"vendor"})
	require.NoError(t, err)

	// vendor files should not be counted
	for _, mod := range ctx.Modules {
		assert.NotContains(t, mod.Dir, "vendor", "vendor dir should be excluded from modules")
	}

	// Files count should only include internal/bar/bar.go (not vendor)
	assert.Equal(t, 1, ctx.Files[".go"], "only 1 .go file should be counted (vendor excluded)")
}

func TestBuildScanContext_SkipHiddenDirs(t *testing.T) {
	root := t.TempDir()

	// .git/ with a go file (should be skipped)
	writeFile(t, filepath.Join(root, ".git", "hooks", "hook.go"), "package hooks\n")
	// src/ with a real go file
	writeFile(t, filepath.Join(root, "go.mod"), "module example.com/myapp\n\ngo 1.21\n")
	writeFile(t, filepath.Join(root, "src", "app.go"), "package src\n")

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	assert.Equal(t, 1, ctx.Files[".go"], "only src/app.go should be counted (.git skipped)")
}

func TestBuildScanContext_EmptyProject(t *testing.T) {
	root := t.TempDir()

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	assert.Equal(t, "unknown", ctx.PrimaryLanguage)
	assert.Equal(t, 0, ctx.TotalFiles)
	assert.NotNil(t, ctx.Files)
	assert.NotNil(t, ctx.Modules, "Modules must not be nil even for empty project")
	assert.NotNil(t, ctx.ExcludedDirs, "ExcludedDirs must not be nil")
	assert.Len(t, ctx.Modules, 0)
}

func TestBuildScanContext_ExistingSpecs(t *testing.T) {
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "go.mod"), "module example.com/myapp\n\ngo 1.21\n")

	// Simulate openspec/changes/{name}/ spec dirs
	require.NoError(t, os.MkdirAll(filepath.Join(root, "openspec", "specs"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "openspec", "changes", "auth-feature"), 0755))
	writeFile(t, filepath.Join(root, "openspec", "changes", "auth-feature", "proposal.md"), "# Auth Feature\n")

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	assert.Contains(t, ctx.ExistingSpecs, "auth-feature", "ExistingSpecs should contain auth-feature")
}
