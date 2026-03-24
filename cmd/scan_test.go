package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mysd/internal/scanner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunScanContextOnly_ValidJSON(t *testing.T) {
	root := t.TempDir()

	// Create a simple Go package in the temp dir
	pkgDir := filepath.Join(root, "mypkg")
	require.NoError(t, os.MkdirAll(pkgDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "mypkg.go"), []byte("package mypkg\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "mypkg_test.go"), []byte("package mypkg_test\n"), 0644))

	var buf bytes.Buffer
	err := runScanContextOnly(&buf, root, nil)
	require.NoError(t, err)

	var ctx scanner.ScanContext
	require.NoError(t, json.Unmarshal(buf.Bytes(), &ctx), "output must be valid JSON")

	assert.Equal(t, root, ctx.RootDir)
	assert.Len(t, ctx.Packages, 1)
	assert.Equal(t, "mypkg", ctx.Packages[0].Name)
	assert.Contains(t, ctx.Packages[0].GoFiles, "mypkg.go")
	assert.Contains(t, ctx.Packages[0].TestFiles, "mypkg_test.go")
}

func TestRunScanContextOnly_EmptyExclude(t *testing.T) {
	root := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(root, "vendor", "lib"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "vendor", "lib", "lib.go"), []byte("package lib\n"), 0644))

	var buf bytes.Buffer
	err := runScanContextOnly(&buf, root, []string{"vendor"})
	require.NoError(t, err)

	var ctx scanner.ScanContext
	require.NoError(t, json.Unmarshal(buf.Bytes(), &ctx))
	assert.Empty(t, ctx.Packages, "vendor package should be excluded")
}

func TestScanCmd_NoFlagsReturnsError(t *testing.T) {
	// Test that running scan without --context-only returns a usage error
	cmd := scanCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runScan(cmd, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "usage: mysd scan --context-only")
}

func TestSetVersion(t *testing.T) {
	// Verify SetVersion sets rootCmd.Version
	SetVersion("v1.2.3-test")
	assert.Equal(t, "v1.2.3-test", rootCmd.Version)

	// Reset to avoid polluting other tests
	SetVersion("dev")
}
