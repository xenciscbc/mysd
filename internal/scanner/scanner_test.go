package scanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mysd/internal/scanner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeFile creates a file with given content in a temp directory.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
}

func TestBuildScanContext_BasicGoProject(t *testing.T) {
	root := t.TempDir()

	// cmd/ package with main.go
	writeFile(t, filepath.Join(root, "cmd", "main.go"), "package main\n")
	// internal/foo/ package with foo.go and foo_test.go
	writeFile(t, filepath.Join(root, "internal", "foo", "foo.go"), "package foo\n")
	writeFile(t, filepath.Join(root, "internal", "foo", "foo_test.go"), "package foo_test\n")

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	assert.Equal(t, root, ctx.RootDir)
	assert.Len(t, ctx.Packages, 2, "expected 2 packages: cmd and internal/foo")
	assert.GreaterOrEqual(t, ctx.TotalFiles, 3)
}

func TestBuildScanContext_ExcludeDirs(t *testing.T) {
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "vendor", "lib", "lib.go"), "package lib\n")
	writeFile(t, filepath.Join(root, "internal", "bar", "bar.go"), "package bar\n")

	ctx, err := scanner.BuildScanContext(root, []string{"vendor"})
	require.NoError(t, err)

	for _, pkg := range ctx.Packages {
		assert.NotContains(t, pkg.Dir, "vendor", "vendor dir should be excluded")
	}
	assert.Len(t, ctx.Packages, 1, "only internal/bar should remain")
}

func TestBuildScanContext_SkipHiddenDirs(t *testing.T) {
	root := t.TempDir()

	// .git/ directory with a go file (should be skipped)
	writeFile(t, filepath.Join(root, ".git", "hooks", "hook.go"), "package hooks\n")
	// src/ with a real go file
	writeFile(t, filepath.Join(root, "src", "app.go"), "package src\n")

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	assert.Len(t, ctx.Packages, 1, "only src package expected")
	assert.Equal(t, "src", ctx.Packages[0].Name)
}

func TestBuildScanContext_ExistingSpecsDetected(t *testing.T) {
	root := t.TempDir()

	// Go package "auth"
	writeFile(t, filepath.Join(root, "auth", "auth.go"), "package auth\n")

	// .specs/changes/auth directory to simulate existing spec
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".specs", "changes", "auth"), 0755))

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	require.Len(t, ctx.Packages, 1)
	assert.True(t, ctx.Packages[0].HasSpec, "auth package should have HasSpec=true")
	assert.Contains(t, ctx.ExistingSpecs, "auth")
}

func TestBuildScanContext_EmptyProject(t *testing.T) {
	root := t.TempDir()

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	assert.Empty(t, ctx.Packages)
	assert.Equal(t, 0, ctx.TotalFiles)
}

func TestBuildScanContext_TestFilesTracked(t *testing.T) {
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "mypkg", "foo.go"), "package mypkg\n")
	writeFile(t, filepath.Join(root, "mypkg", "foo_test.go"), "package mypkg_test\n")

	ctx, err := scanner.BuildScanContext(root, nil)
	require.NoError(t, err)

	require.Len(t, ctx.Packages, 1)
	pkg := ctx.Packages[0]
	assert.Contains(t, pkg.GoFiles, "foo.go")
	assert.Contains(t, pkg.TestFiles, "foo_test.go")
}
