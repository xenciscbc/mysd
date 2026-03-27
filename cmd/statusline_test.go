package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// helperStatuslineCmd creates an isolated cobra.Command for testing statusline.
func helperStatuslineCmd(outBuf *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{Use: "statusline"}
	cmd.SetOut(outBuf)
	return cmd
}

// helperReadStatuslineEnabled reads statusline_enabled from a mysd.yaml in the given dir.
// Returns (value, exists, error).
func helperReadStatuslineEnabled(dir string) (bool, bool, error) {
	data, err := os.ReadFile(filepath.Join(dir, ".claude", "mysd.yaml"))
	if err != nil {
		return false, false, err
	}
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return false, false, err
	}
	val, ok := raw["statusline_enabled"]
	if !ok {
		return false, false, nil
	}
	b, _ := val.(bool)
	return b, true, nil
}

// helperWriteStatuslineEnabled pre-writes a mysd.yaml with a known statusline_enabled value.
func helperWriteStatuslineEnabled(t *testing.T, dir string, enabled bool) {
	t.Helper()
	claudeDir := filepath.Join(dir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := map[string]interface{}{"statusline_enabled": enabled}
	data, err := yaml.Marshal(content)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), data, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestRunStatuslineOn(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	cmd := helperStatuslineCmd(&buf)

	if err := runStatuslineInDir(cmd, []string{"on"}, dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, exists, err := helperReadStatuslineEnabled(dir)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if !exists {
		t.Fatal("statusline_enabled not written to mysd.yaml")
	}
	if !val {
		t.Errorf("expected statusline_enabled=true, got false")
	}
}

func TestRunStatuslineOff(t *testing.T) {
	dir := t.TempDir()
	// Pre-set to true so we can verify it flips to false.
	helperWriteStatuslineEnabled(t, dir, true)

	var buf bytes.Buffer
	cmd := helperStatuslineCmd(&buf)

	if err := runStatuslineInDir(cmd, []string{"off"}, dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, exists, err := helperReadStatuslineEnabled(dir)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if !exists {
		t.Fatal("statusline_enabled not written to mysd.yaml")
	}
	if val {
		t.Errorf("expected statusline_enabled=false, got true")
	}
}

func TestRunStatuslineToggle(t *testing.T) {
	// Toggle from true -> false
	dir := t.TempDir()
	helperWriteStatuslineEnabled(t, dir, true)

	var buf bytes.Buffer
	cmd := helperStatuslineCmd(&buf)

	if err := runStatuslineInDir(cmd, []string{}, dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, exists, err := helperReadStatuslineEnabled(dir)
	if err != nil || !exists {
		t.Fatalf("failed to read config: err=%v exists=%v", err, exists)
	}
	if val {
		t.Errorf("expected toggle from true -> false, got true")
	}

	// Toggle from false -> true
	dir2 := t.TempDir()
	helperWriteStatuslineEnabled(t, dir2, false)

	var buf2 bytes.Buffer
	cmd2 := helperStatuslineCmd(&buf2)

	if err := runStatuslineInDir(cmd2, []string{}, dir2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val2, exists2, err2 := helperReadStatuslineEnabled(dir2)
	if err2 != nil || !exists2 {
		t.Fatalf("failed to read config: err=%v exists=%v", err2, exists2)
	}
	if !val2 {
		t.Errorf("expected toggle from false -> true, got false")
	}
}

func TestRunStatuslineToggleDefault(t *testing.T) {
	// No existing field -> default is true, toggle -> false
	dir := t.TempDir()
	// Create empty .claude/mysd.yaml (no statusline_enabled field)
	claudeDir := filepath.Join(dir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "mysd.yaml"), []byte("model_profile: balanced\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	cmd := helperStatuslineCmd(&buf)

	if err := runStatuslineInDir(cmd, []string{}, dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, exists, err := helperReadStatuslineEnabled(dir)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if !exists {
		t.Fatal("statusline_enabled not written after toggle from default")
	}
	if val {
		t.Errorf("expected toggle from default(true) -> false, got true")
	}
}

func TestStatuslineOutputFormat(t *testing.T) {
	tests := []struct {
		arg      string
		expected string
	}{
		{"on", "Statusline: on\n"},
		{"off", "Statusline: off\n"},
	}

	for _, tc := range tests {
		dir := t.TempDir()
		var buf bytes.Buffer
		cmd := helperStatuslineCmd(&buf)

		if err := runStatuslineInDir(cmd, []string{tc.arg}, dir); err != nil {
			t.Fatalf("arg=%q unexpected error: %v", tc.arg, err)
		}

		if got := buf.String(); got != tc.expected {
			t.Errorf("arg=%q: expected output %q, got %q", tc.arg, tc.expected, got)
		}
	}
}
