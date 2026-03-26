package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultOwner is the GitHub repository owner.
	DefaultOwner = "xenciscbc"
	// DefaultRepo is the GitHub repository name.
	DefaultRepo = "mysd"
	// GitHubAPIBase is the base URL for GitHub API.
	GitHubAPIBase = "https://api.github.com"
)

// Semver represents a semantic version with Major, Minor, and Patch components.
type Semver struct {
	Major, Minor, Patch int
}

// ReleaseInfo represents the GitHub Releases API response for the latest release.
type ReleaseInfo struct {
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Assets  []Asset `json:"assets"`
	HTMLURL string  `json:"html_url"`
}

// Asset represents a single release asset in the GitHub API response.
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// ParseSemver parses a semantic version string (with or without "v" prefix).
// Returns an error for "dev", empty strings, or malformed versions.
func ParseSemver(s string) (Semver, error) {
	// Strip leading "v" prefix
	s = strings.TrimPrefix(s, "v")

	if s == "" || s == "dev" {
		return Semver{}, fmt.Errorf("update: cannot parse version %q as semver", s)
	}

	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return Semver{}, fmt.Errorf("update: version %q must have 3 components (major.minor.patch)", s)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Semver{}, fmt.Errorf("update: invalid major version in %q: %w", s, err)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Semver{}, fmt.Errorf("update: invalid minor version in %q: %w", s, err)
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return Semver{}, fmt.Errorf("update: invalid patch version in %q: %w", s, err)
	}

	return Semver{Major: major, Minor: minor, Patch: patch}, nil
}

// LessThan returns true if s is an older version than other.
// Comparison order: Major > Minor > Patch.
func (s Semver) LessThan(other Semver) bool {
	if s.Major != other.Major {
		return s.Major < other.Major
	}
	if s.Minor != other.Minor {
		return s.Minor < other.Minor
	}
	return s.Patch < other.Patch
}

// IsUpdateAvailable returns true if currentVersion is older than latestVersion.
// Per D-02, "dev" is always considered outdated compared to any release version.
func IsUpdateAvailable(currentVersion, latestVersion string) (bool, error) {
	// dev build is always outdated
	if currentVersion == "dev" {
		return true, nil
	}

	current, err := ParseSemver(currentVersion)
	if err != nil {
		return false, fmt.Errorf("update: invalid current version: %w", err)
	}

	latest, err := ParseSemver(latestVersion)
	if err != nil {
		return false, fmt.Errorf("update: invalid latest version: %w", err)
	}

	return current.LessThan(latest), nil
}

// CheckLatestVersion queries the GitHub Releases API for the latest release.
// If httpClient is nil, a default client with 15-second timeout is used.
// Uses DefaultOwner and DefaultRepo as the GitHub repository coordinates.
func CheckLatestVersion(ctx context.Context, httpClient *http.Client) (ReleaseInfo, error) {
	return CheckLatestVersionWithBase(ctx, httpClient, GitHubAPIBase)
}

// CheckLatestVersionWithBase is the testable version of CheckLatestVersion that
// accepts a custom API base URL (for testing with httptest.NewServer).
func CheckLatestVersionWithBase(ctx context.Context, httpClient *http.Client, apiBase string) (ReleaseInfo, error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}

	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", apiBase, DefaultOwner, DefaultRepo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ReleaseInfo{}, fmt.Errorf("update: failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return ReleaseInfo{}, fmt.Errorf("update: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ReleaseInfo{}, fmt.Errorf("update: GitHub API returned status %d", resp.StatusCode)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return ReleaseInfo{}, fmt.Errorf("update: failed to decode response: %w", err)
	}

	return release, nil
}

// AssetNameForPlatform returns the GoReleaser-generated archive name for the given
// platform. Matches name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}".
// Returns .zip for windows, .tar.gz for all other platforms.
func AssetNameForPlatform(goos, goarch, version string) string {
	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}
	return fmt.Sprintf("mysd_%s_%s_%s%s", version, goos, goarch, ext)
}

// FindAssetURL searches release.Assets for the given asset name and returns
// its BrowserDownloadURL. Returns an error if the asset is not found.
func FindAssetURL(release ReleaseInfo, assetName string) (string, error) {
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return asset.BrowserDownloadURL, nil
		}
	}
	return "", fmt.Errorf("update: asset %q not found in release %s", assetName, release.TagName)
}

// FindChecksumURL searches release.Assets for "checksums.txt" and returns its URL.
// Returns an error if checksums.txt is not found.
func FindChecksumURL(release ReleaseInfo) (string, error) {
	return FindAssetURL(release, "checksums.txt")
}
