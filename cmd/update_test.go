package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUpdateCmdRegistered verifies updateCmd is registered on rootCmd.
func TestUpdateCmdRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "update" {
			return
		}
	}
	t.Fatal("updateCmd not registered on rootCmd")
}

// TestUpdateCmdFlags verifies --check and --force flags exist on updateCmd.
func TestUpdateCmdFlags(t *testing.T) {
	checkFlag := updateCmd.Flags().Lookup("check")
	require.NotNil(t, checkFlag, "--check flag must be defined")
	assert.Equal(t, "bool", checkFlag.Value.Type())

	forceFlag := updateCmd.Flags().Lookup("force")
	require.NotNil(t, forceFlag, "--force flag must be defined")
	assert.Equal(t, "bool", forceFlag.Value.Type())
}

// TestUpdateCmdUsage verifies the Use field and Short description.
func TestUpdateCmdUsage(t *testing.T) {
	assert.Equal(t, "update", updateCmd.Use)
	assert.NotEmpty(t, updateCmd.Short)
}

// TestUpdateOutputStruct verifies UpdateOutput marshals to expected JSON keys.
func TestUpdateOutputStruct(t *testing.T) {
	out := UpdateOutput{
		CurrentVersion:  "v1.0.0",
		LatestVersion:   "v1.1.0",
		UpdateAvailable: true,
		ReleaseURL:      "https://github.com/xenciscbc/mysd/releases/tag/v1.1.0",
		CheckOnly:       true,
		Force:           false,
		BinaryUpdated:   false,
		PluginSync:      nil,
		Error:           "",
	}

	data, err := json.MarshalIndent(out, "", "  ")
	require.NoError(t, err)

	var decoded map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, "v1.0.0", decoded["current_version"])
	assert.Equal(t, "v1.1.0", decoded["latest_version"])
	assert.Equal(t, true, decoded["update_available"])
	assert.Equal(t, true, decoded["check_only"])
}

// TestUpdateOutputOmitsEmptyFields verifies omitempty fields are absent when empty.
func TestUpdateOutputOmitsEmptyFields(t *testing.T) {
	out := UpdateOutput{
		CurrentVersion:  "dev",
		UpdateAvailable: false,
		CheckOnly:       false,
		Force:           false,
		BinaryUpdated:   false,
	}

	data, err := json.MarshalIndent(out, "", "  ")
	require.NoError(t, err)

	s := string(data)
	// omitempty fields should not appear when empty
	assert.False(t, strings.Contains(s, `"latest_version"`), "latest_version should be omitted when empty")
	assert.False(t, strings.Contains(s, `"release_url"`), "release_url should be omitted when empty")
	assert.False(t, strings.Contains(s, `"plugin_sync"`), "plugin_sync should be omitted when nil")
	assert.False(t, strings.Contains(s, `"error"`), "error should be omitted when empty")
}

// TestSyncOutputStruct verifies SyncOutput marshals correctly.
func TestSyncOutputStruct(t *testing.T) {
	so := SyncOutput{
		Added:   2,
		Updated: 1,
		Deleted: 0,
		Errors:  []string{"delete command foo.md: file not found"},
	}

	data, err := json.MarshalIndent(so, "", "  ")
	require.NoError(t, err)

	var decoded map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, float64(2), decoded["added"])
	assert.Equal(t, float64(1), decoded["updated"])
}

// TestPrintJSON verifies printJSON writes indented JSON to stdout.
func TestPrintJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	cmd := updateCmd
	cmd.SetOut(buf)

	err := printJSON(cmd, map[string]string{"key": "value"})
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, `"key"`)
	assert.Contains(t, output, `"value"`)
}
