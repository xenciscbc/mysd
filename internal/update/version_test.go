package update_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xenciscbc/mysd/internal/update"
)

func TestParseSemver(t *testing.T) {
	t.Run("valid version without v prefix", func(t *testing.T) {
		sv, err := update.ParseSemver("1.2.3")
		require.NoError(t, err)
		assert.Equal(t, 1, sv.Major)
		assert.Equal(t, 2, sv.Minor)
		assert.Equal(t, 3, sv.Patch)
	})

	t.Run("valid version with v prefix", func(t *testing.T) {
		sv, err := update.ParseSemver("v1.2.3")
		require.NoError(t, err)
		assert.Equal(t, 1, sv.Major)
		assert.Equal(t, 2, sv.Minor)
		assert.Equal(t, 3, sv.Patch)
	})

	t.Run("dev version returns error", func(t *testing.T) {
		_, err := update.ParseSemver("dev")
		require.Error(t, err)
	})

	t.Run("invalid version returns error", func(t *testing.T) {
		_, err := update.ParseSemver("invalid")
		require.Error(t, err)
	})

	t.Run("partial version returns error", func(t *testing.T) {
		_, err := update.ParseSemver("1.2")
		require.Error(t, err)
	})
}

func TestIsUpdateAvailable(t *testing.T) {
	t.Run("dev version is always outdated", func(t *testing.T) {
		available, err := update.IsUpdateAvailable("dev", "1.0.0")
		require.NoError(t, err)
		assert.True(t, available)
	})

	t.Run("older current version needs update", func(t *testing.T) {
		available, err := update.IsUpdateAvailable("1.0.0", "1.0.1")
		require.NoError(t, err)
		assert.True(t, available)
	})

	t.Run("newer current version does not need update", func(t *testing.T) {
		available, err := update.IsUpdateAvailable("1.0.1", "1.0.0")
		require.NoError(t, err)
		assert.False(t, available)
	})

	t.Run("same version does not need update", func(t *testing.T) {
		available, err := update.IsUpdateAvailable("1.0.0", "1.0.0")
		require.NoError(t, err)
		assert.False(t, available)
	})

	t.Run("minor version increment needs update", func(t *testing.T) {
		available, err := update.IsUpdateAvailable("1.0.0", "1.1.0")
		require.NoError(t, err)
		assert.True(t, available)
	})

	t.Run("major version increment needs update", func(t *testing.T) {
		available, err := update.IsUpdateAvailable("1.0.0", "2.0.0")
		require.NoError(t, err)
		assert.True(t, available)
	})
}

func TestCheckLatestVersion(t *testing.T) {
	t.Run("valid GitHub API response", func(t *testing.T) {
		release := update.ReleaseInfo{
			TagName: "v1.2.3",
			Name:    "Release 1.2.3",
			HTMLURL: "https://github.com/xenciscbc/mysd/releases/tag/v1.2.3",
			Assets: []update.Asset{
				{
					Name:               "mysd_1.2.3_linux_amd64.tar.gz",
					BrowserDownloadURL: "https://github.com/xenciscbc/mysd/releases/download/v1.2.3/mysd_1.2.3_linux_amd64.tar.gz",
				},
				{
					Name:               "checksums.txt",
					BrowserDownloadURL: "https://github.com/xenciscbc/mysd/releases/download/v1.2.3/checksums.txt",
				},
			},
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "/releases/latest")
			assert.Equal(t, "application/vnd.github+json", r.Header.Get("Accept"))
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(release)
		}))
		defer srv.Close()

		client := &http.Client{Timeout: 5 * time.Second}
		info, err := update.CheckLatestVersionWithBase(context.Background(), client, srv.URL)
		require.NoError(t, err)
		assert.Equal(t, "v1.2.3", info.TagName)
		assert.Len(t, info.Assets, 2)
	})

	t.Run("404 response returns error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		client := &http.Client{Timeout: 5 * time.Second}
		_, err := update.CheckLatestVersionWithBase(context.Background(), client, srv.URL)
		require.Error(t, err)
	})

	t.Run("network timeout returns error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate slow server by sleeping longer than client timeout
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		client := &http.Client{Timeout: 10 * time.Millisecond}
		_, err := update.CheckLatestVersionWithBase(context.Background(), client, srv.URL)
		require.Error(t, err)
	})
}

func TestAssetNameForPlatform(t *testing.T) {
	t.Run("linux amd64", func(t *testing.T) {
		name := update.AssetNameForPlatform("linux", "amd64", "1.0.0")
		assert.Equal(t, "mysd_1.0.0_linux_amd64.tar.gz", name)
	})

	t.Run("windows amd64", func(t *testing.T) {
		name := update.AssetNameForPlatform("windows", "amd64", "1.0.0")
		assert.Equal(t, "mysd_1.0.0_windows_amd64.zip", name)
	})

	t.Run("darwin arm64", func(t *testing.T) {
		name := update.AssetNameForPlatform("darwin", "arm64", "1.0.0")
		assert.Equal(t, "mysd_1.0.0_darwin_arm64.tar.gz", name)
	})
}

func TestFindAssetURL(t *testing.T) {
	release := update.ReleaseInfo{
		Assets: []update.Asset{
			{
				Name:               "mysd_1.0.0_linux_amd64.tar.gz",
				BrowserDownloadURL: "https://example.com/mysd_1.0.0_linux_amd64.tar.gz",
			},
			{
				Name:               "checksums.txt",
				BrowserDownloadURL: "https://example.com/checksums.txt",
			},
		},
	}

	t.Run("finds existing asset", func(t *testing.T) {
		url, err := update.FindAssetURL(release, "mysd_1.0.0_linux_amd64.tar.gz")
		require.NoError(t, err)
		assert.Equal(t, "https://example.com/mysd_1.0.0_linux_amd64.tar.gz", url)
	})

	t.Run("returns error for missing asset", func(t *testing.T) {
		_, err := update.FindAssetURL(release, "nonexistent.tar.gz")
		require.Error(t, err)
	})
}

func TestFindChecksumURL(t *testing.T) {
	t.Run("finds checksums.txt", func(t *testing.T) {
		release := update.ReleaseInfo{
			Assets: []update.Asset{
				{
					Name:               "checksums.txt",
					BrowserDownloadURL: "https://example.com/checksums.txt",
				},
			},
		}
		url, err := update.FindChecksumURL(release)
		require.NoError(t, err)
		assert.Equal(t, "https://example.com/checksums.txt", url)
	})

	t.Run("returns error when no checksums.txt", func(t *testing.T) {
		release := update.ReleaseInfo{
			Assets: []update.Asset{},
		}
		_, err := update.FindChecksumURL(release)
		require.Error(t, err)
	})
}
