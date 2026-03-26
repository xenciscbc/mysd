package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xenciscbc/mysd/internal/scanner"
)

func TestRunScanContextOnly_ValidJSON(t *testing.T) {
	root := t.TempDir()

	// Create a simple Go package in the temp dir
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/mypkg\n\ngo 1.21\n"), 0644))
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
	assert.Equal(t, "go", ctx.PrimaryLanguage)
	assert.Contains(t, ctx.Files, ".go", "Files map should contain .go key")
	assert.GreaterOrEqual(t, ctx.Files[".go"], 2, "should have at least 2 .go files")
}

func TestRunScanContextOnly_ExcludeVendor(t *testing.T) {
	root := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/myapp\n\ngo 1.21\n"), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "vendor", "lib"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "vendor", "lib", "lib.go"), []byte("package lib\n"), 0644))

	var buf bytes.Buffer
	err := runScanContextOnly(&buf, root, []string{"vendor"})
	require.NoError(t, err)

	var ctx scanner.ScanContext
	require.NoError(t, json.Unmarshal(buf.Bytes(), &ctx))
	// go.mod is not a .go file, so Files[".go"] should be 0 since vendor is excluded
	assert.Equal(t, 0, ctx.Files[".go"], "vendor .go files should be excluded")
}

func TestScanCmd_NoFlagsReturnsError(t *testing.T) {
	// Test that running scan without --context-only or --scaffold-only returns a usage error
	cmd := scanCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runScan(cmd, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "usage: mysd scan --context-only|--scaffold-only")
}

func TestScaffoldOpenSpecDir(t *testing.T) {
	root := t.TempDir()

	err := scaffoldOpenSpecDir(root)
	require.NoError(t, err)

	// Check that openspec/ and openspec/specs/ were created
	info, err := os.Stat(filepath.Join(root, "openspec"))
	require.NoError(t, err)
	assert.True(t, info.IsDir(), "openspec/ should be a directory")

	info, err = os.Stat(filepath.Join(root, "openspec", "specs"))
	require.NoError(t, err)
	assert.True(t, info.IsDir(), "openspec/specs/ should be a directory")

	// Should NOT create openspec/config.yaml (per D-06)
	_, err = os.Stat(filepath.Join(root, "openspec", "config.yaml"))
	assert.True(t, os.IsNotExist(err), "openspec/config.yaml should NOT be created by scaffoldOpenSpecDir")
}

func TestScaffoldOpenSpecDir_Idempotent(t *testing.T) {
	root := t.TempDir()

	// Call twice — should not error
	require.NoError(t, scaffoldOpenSpecDir(root))
	require.NoError(t, scaffoldOpenSpecDir(root))
}

func TestSetVersion(t *testing.T) {
	// Verify SetVersion sets rootCmd.Version
	SetVersion("v1.2.3-test")
	assert.Equal(t, "v1.2.3-test", rootCmd.Version)

	// Reset to avoid polluting other tests
	SetVersion("dev")
}
