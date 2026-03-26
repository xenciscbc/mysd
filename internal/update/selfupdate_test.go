package update_test

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xenciscbc/mysd/internal/update"
)

// createTestFile creates a temp file with given content and returns its path.
func createTestFile(t *testing.T, content []byte) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*")
	require.NoError(t, err)
	_, err = f.Write(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return f.Name()
}

// sha256Hex returns the hex-encoded SHA256 hash of content.
func sha256Hex(content []byte) string {
	h := sha256.Sum256(content)
	return hex.EncodeToString(h[:])
}

func TestVerifyChecksum(t *testing.T) {
	content := []byte("test binary content for checksum verification")
	expectedHash := sha256Hex(content)

	filePath := createTestFile(t, content)

	t.Run("matching hash passes", func(t *testing.T) {
		err := update.VerifyChecksum(filePath, expectedHash)
		require.NoError(t, err)
	})

	t.Run("mismatching hash returns error", func(t *testing.T) {
		err := update.VerifyChecksum(filePath, "0000000000000000000000000000000000000000000000000000000000000000")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "checksum mismatch")
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		err := update.VerifyChecksum("/nonexistent/path/file.bin", expectedHash)
		require.Error(t, err)
	})
}

func TestParseChecksumFile(t *testing.T) {
	// GoReleaser checksums.txt format: "{hash}  {filename}\n" (double space)
	checksumContent := `abc123def456abc123def456abc123def456abc123def456abc123def456abc1  mysd_1.0.0_linux_amd64.tar.gz
def456abc123def456abc123def456abc123def456abc123def456abc123def4  mysd_1.0.0_darwin_arm64.tar.gz
ffff0000ffff0000ffff0000ffff0000ffff0000ffff0000ffff0000ffff0000  mysd_1.0.0_windows_amd64.zip
`

	t.Run("finds correct hash for existing filename", func(t *testing.T) {
		hash, err := update.ParseChecksumFile(checksumContent, "mysd_1.0.0_linux_amd64.tar.gz")
		require.NoError(t, err)
		assert.Equal(t, "abc123def456abc123def456abc123def456abc123def456abc123def456abc1", hash)
	})

	t.Run("finds correct hash for another filename", func(t *testing.T) {
		hash, err := update.ParseChecksumFile(checksumContent, "mysd_1.0.0_windows_amd64.zip")
		require.NoError(t, err)
		assert.Equal(t, "ffff0000ffff0000ffff0000ffff0000ffff0000ffff0000ffff0000ffff0000", hash)
	})

	t.Run("returns error for filename not found", func(t *testing.T) {
		_, err := update.ParseChecksumFile(checksumContent, "nonexistent_file.tar.gz")
		require.Error(t, err)
	})
}

func TestExtractBinary(t *testing.T) {
	binaryContent := []byte("#!/bin/sh\necho 'hello from mysd'")

	t.Run("extract from tar.gz", func(t *testing.T) {
		// Create a test tar.gz archive containing "mysd" binary
		archivePath := createTarGzWithBinary(t, "mysd", binaryContent)

		extractedPath, err := update.ExtractBinary(archivePath, false)
		require.NoError(t, err)
		defer os.Remove(extractedPath)

		content, err := os.ReadFile(extractedPath)
		require.NoError(t, err)
		assert.Equal(t, binaryContent, content)
	})

	t.Run("extract from zip", func(t *testing.T) {
		// Create a test zip archive containing "mysd.exe" binary
		archivePath := createZipWithBinary(t, "mysd.exe", binaryContent)

		extractedPath, err := update.ExtractBinary(archivePath, true)
		require.NoError(t, err)
		defer os.Remove(extractedPath)

		content, err := os.ReadFile(extractedPath)
		require.NoError(t, err)
		assert.Equal(t, binaryContent, content)
	})
}

func TestRollback(t *testing.T) {
	t.Run("restores .old file to original path", func(t *testing.T) {
		dir := t.TempDir()
		exePath := filepath.Join(dir, "mysd")
		oldContent := []byte("old binary content")

		// Create .old file (simulating what replaceExecutable does)
		oldPath := exePath + ".old"
		require.NoError(t, os.WriteFile(oldPath, oldContent, 0755))

		err := update.Rollback(exePath)
		require.NoError(t, err)

		// Verify .old was restored to original path
		content, err := os.ReadFile(exePath)
		require.NoError(t, err)
		assert.Equal(t, oldContent, content)

		// Verify .old file no longer exists
		_, err = os.Stat(oldPath)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("returns error when .old file does not exist", func(t *testing.T) {
		dir := t.TempDir()
		exePath := filepath.Join(dir, "mysd_nonexistent")

		err := update.Rollback(exePath)
		require.Error(t, err)
	})
}

func TestDownloadFile(t *testing.T) {
	content := []byte("downloaded content for testing")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(content)
	}))
	defer srv.Close()

	t.Run("downloads file to temp path", func(t *testing.T) {
		tmpPath, err := update.DownloadFile(t.Context(), srv.URL+"/test-file", nil)
		require.NoError(t, err)
		defer os.Remove(tmpPath)

		downloaded, err := os.ReadFile(tmpPath)
		require.NoError(t, err)
		assert.Equal(t, content, downloaded)
	})
}

// createTarGzWithBinary creates a test .tar.gz archive with a binary file inside.
func createTarGzWithBinary(t *testing.T, binaryName string, content []byte) string {
	t.Helper()

	archivePath := filepath.Join(t.TempDir(), "test-archive.tar.gz")
	f, err := os.Create(archivePath)
	require.NoError(t, err)
	defer f.Close()

	gzw := gzip.NewWriter(f)
	tw := tar.NewWriter(gzw)

	hdr := &tar.Header{
		Name: binaryName,
		Mode: 0755,
		Size: int64(len(content)),
	}
	require.NoError(t, tw.WriteHeader(hdr))
	_, err = tw.Write(content)
	require.NoError(t, err)
	require.NoError(t, tw.Close())
	require.NoError(t, gzw.Close())

	return archivePath
}

// createZipWithBinary creates a test .zip archive with a binary file inside.
func createZipWithBinary(t *testing.T, binaryName string, content []byte) string {
	t.Helper()

	archivePath := filepath.Join(t.TempDir(), "test-archive.zip")
	f, err := os.Create(archivePath)
	require.NoError(t, err)
	defer f.Close()

	zw := zip.NewWriter(f)
	fw, err := zw.Create(binaryName)
	require.NoError(t, err)
	_, err = fw.Write(content)
	require.NoError(t, err)
	require.NoError(t, zw.Close())

	return archivePath
}

// TestApplyUpdateChecksumMismatch verifies that ApplyUpdate rejects a bad archive.
func TestApplyUpdateChecksumMismatch(t *testing.T) {
	binaryContent := []byte("fake binary")
	archiveContent := []byte("not actually a tar.gz")
	badHash := "0000000000000000000000000000000000000000000000000000000000000000"

	// Use the current platform's asset name so ApplyUpdate can find it in the release.
	assetName := update.AssetNameForPlatform(runtime.GOOS, runtime.GOARCH, "1.0.0")

	// Serve archive and checksums.txt
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/checksums.txt":
			// Bad hash for the archive — this is what triggers checksum mismatch
			fmt.Fprintf(w, "%s  %s\n", badHash, assetName)
		default:
			_, _ = w.Write(archiveContent)
		}
	}))
	defer srv.Close()

	dir := t.TempDir()
	currentExe := filepath.Join(dir, "mysd")
	require.NoError(t, os.WriteFile(currentExe, binaryContent, 0755))

	release := ReleaseInfo{
		TagName: "v1.0.0",
		Assets: []Asset{
			{
				Name:               assetName,
				BrowserDownloadURL: srv.URL + "/" + assetName,
			},
			{
				Name:               "checksums.txt",
				BrowserDownloadURL: srv.URL + "/checksums.txt",
			},
		},
	}

	err := update.ApplyUpdate(t.Context(), nil, release, currentExe)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "checksum")

	// Verify original binary was not replaced
	content, err := os.ReadFile(currentExe)
	require.NoError(t, err)
	assert.Equal(t, binaryContent, content)
}

// ReleaseInfo and Asset aliases for test package usage
type ReleaseInfo = update.ReleaseInfo
type Asset = update.Asset
