package update

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// DownloadFile downloads the resource at url to a temporary file.
// Returns the path to the downloaded temp file.
// If httpClient is nil, a default client with 30-second timeout is used.
func DownloadFile(ctx context.Context, url string, httpClient *http.Client) (string, error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("update: failed to create download request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("update: download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("update: download returned status %d for %s", resp.StatusCode, url)
	}

	tmpFile, err := os.CreateTemp("", "mysd-download-*")
	if err != nil {
		return "", fmt.Errorf("update: failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("update: failed to write download: %w", err)
	}

	return tmpFile.Name(), nil
}

// VerifyChecksum computes the SHA256 hash of the file at filePath and compares
// it against expectedHex. Returns an error containing "checksum mismatch" if they differ.
func VerifyChecksum(filePath, expectedHex string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("update: cannot open file for checksum: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("update: failed to hash file: %w", err)
	}

	actualHex := hex.EncodeToString(h.Sum(nil))
	if actualHex != strings.ToLower(expectedHex) {
		return fmt.Errorf("update: checksum mismatch for %s: expected %s, got %s", filePath, expectedHex, actualHex)
	}

	return nil
}

// ParseChecksumFile parses a GoReleaser checksums.txt file.
// Format: "{hash}  {filename}\n" — double space separator.
// Returns the SHA256 hash for the given filename, or an error if not found.
func ParseChecksumFile(content, filename string) (string, error) {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Split on double space (GoReleaser format)
		parts := strings.SplitN(line, "  ", 2)
		if len(parts) != 2 {
			continue
		}
		hash := strings.TrimSpace(parts[0])
		name := strings.TrimSpace(parts[1])
		if name == filename {
			return hash, nil
		}
	}
	return "", fmt.Errorf("update: filename %q not found in checksums file", filename)
}

// ExtractBinary extracts the mysd binary from a tar.gz or zip archive.
// isZip=true for Windows .zip archives, false for tar.gz.
// Returns the path to the extracted binary in a temp directory.
func ExtractBinary(archivePath string, isZip bool) (string, error) {
	if isZip {
		return extractFromZip(archivePath)
	}
	return extractFromTarGz(archivePath)
}

// extractFromTarGz extracts the mysd binary from a tar.gz archive.
func extractFromTarGz(archivePath string) (string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", fmt.Errorf("update: cannot open archive: %w", err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return "", fmt.Errorf("update: failed to create gzip reader: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("update: tar read error: %w", err)
		}

		// Find the binary — it is a regular file named "mysd" (no directory prefix needed,
		// but we strip directory components to handle both flat and nested archives)
		baseName := filepath.Base(hdr.Name)
		if hdr.Typeflag == tar.TypeReg && isBinaryName(baseName) {
			return writeTempBinary(tr, baseName)
		}
	}

	return "", fmt.Errorf("update: mysd binary not found in tar.gz archive")
}

// extractFromZip extracts the mysd binary from a zip archive.
func extractFromZip(archivePath string) (string, error) {
	zr, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", fmt.Errorf("update: cannot open zip archive: %w", err)
	}
	defer zr.Close()

	for _, f := range zr.File {
		baseName := filepath.Base(f.Name)
		if !f.FileInfo().IsDir() && isBinaryName(baseName) {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("update: cannot open zip entry: %w", err)
			}
			defer rc.Close()
			return writeTempBinary(rc, baseName)
		}
	}

	return "", fmt.Errorf("update: mysd binary not found in zip archive")
}

// isBinaryName returns true if name matches the mysd binary (with or without .exe).
func isBinaryName(name string) bool {
	return name == "mysd" || name == "mysd.exe"
}

// writeTempBinary writes the reader content to a new temp file and returns its path.
func writeTempBinary(r io.Reader, name string) (string, error) {
	tmpFile, err := os.CreateTemp("", "mysd-extract-"+name+"-*")
	if err != nil {
		return "", fmt.Errorf("update: failed to create temp file for extraction: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, r); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("update: failed to write extracted binary: %w", err)
	}

	return tmpFile.Name(), nil
}

// ApplyUpdate orchestrates the full binary self-update:
//  1. Download the release archive
//  2. Download checksums.txt
//  3. Verify archive checksum
//  4. Extract binary from archive
//  5. Call replaceExecutable (platform-specific)
//  6. On any failure after step 5, call Rollback
//
// currentExePath should be the path to the currently running binary (os.Executable()).
func ApplyUpdate(ctx context.Context, httpClient *http.Client, release ReleaseInfo, currentExePath string) error {
	// Determine archive name for current platform
	version := strings.TrimPrefix(release.TagName, "v")
	isZip := runtime.GOOS == "windows"
	assetName := AssetNameForPlatform(runtime.GOOS, runtime.GOARCH, version)

	// Find asset and checksum URLs
	assetURL, err := FindAssetURL(release, assetName)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	checksumURL, err := FindChecksumURL(release)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	// Step 1: Download archive
	archivePath, err := DownloadFile(ctx, assetURL, httpClient)
	if err != nil {
		return fmt.Errorf("update: failed to download archive: %w", err)
	}
	defer os.Remove(archivePath)

	// Step 2: Download checksums.txt
	checksumPath, err := DownloadFile(ctx, checksumURL, httpClient)
	if err != nil {
		return fmt.Errorf("update: failed to download checksums: %w", err)
	}
	defer os.Remove(checksumPath)

	// Step 3: Verify archive checksum
	checksumContent, err := os.ReadFile(checksumPath)
	if err != nil {
		return fmt.Errorf("update: failed to read checksums file: %w", err)
	}
	expectedHash, err := ParseChecksumFile(string(checksumContent), assetName)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	if err := VerifyChecksum(archivePath, expectedHash); err != nil {
		return err
	}

	// Step 4: Extract binary from archive
	newBinaryPath, err := ExtractBinary(archivePath, isZip)
	if err != nil {
		return fmt.Errorf("update: failed to extract binary: %w", err)
	}
	defer os.Remove(newBinaryPath)

	// Step 5: Replace executable (platform-specific)
	if err := replaceExecutable(currentExePath, newBinaryPath); err != nil {
		return fmt.Errorf("update: failed to replace executable: %w", err)
	}

	return nil
}

// Rollback restores {exePath}.old back to exePath.
// Called when an update fails after replaceExecutable has already run.
func Rollback(exePath string) error {
	oldPath := exePath + ".old"
	if err := os.Rename(oldPath, exePath); err != nil {
		return fmt.Errorf("update: rollback failed — cannot restore %s from %s: %w", exePath, oldPath, err)
	}
	return nil
}
