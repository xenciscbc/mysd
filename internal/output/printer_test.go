package output

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewPrinter_NonTTY verifies that bytes.Buffer is detected as non-TTY.
func TestNewPrinter_NonTTY(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	assert.False(t, p.isTTY, "bytes.Buffer should be detected as non-TTY")
}

func TestSuccess_NonTTY_ContainsPrefix(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Success("operation completed")
	output := buf.String()
	assert.True(t, strings.HasPrefix(output, "[OK]"), "non-TTY Success should start with [OK], got: %q", output)
	assert.Contains(t, output, "operation completed")
}

func TestError_NonTTY_ContainsPrefix(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Error("something went wrong")
	output := buf.String()
	assert.True(t, strings.HasPrefix(output, "[ERROR]"), "non-TTY Error should start with [ERROR], got: %q", output)
	assert.Contains(t, output, "something went wrong")
}

func TestWarning_NonTTY_ContainsPrefix(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Warning("check your config")
	output := buf.String()
	assert.True(t, strings.HasPrefix(output, "[WARN]"), "non-TTY Warning should start with [WARN], got: %q", output)
	assert.Contains(t, output, "check your config")
}

func TestInfo_NonTTY_ContainsPrefix(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Info("loading configuration")
	output := buf.String()
	assert.True(t, strings.HasPrefix(output, "[INFO]"), "non-TTY Info should start with [INFO], got: %q", output)
	assert.Contains(t, output, "loading configuration")
}

func TestHeader_NonTTY_ContainsTripleEquals(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Header("My Section")
	output := buf.String()
	assert.Contains(t, output, "===", "non-TTY Header should contain ===")
	assert.Contains(t, output, "My Section")
}

func TestMuted_NonTTY_PlainText(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Muted("some muted text")
	output := buf.String()
	assert.Contains(t, output, "some muted text")
}

func TestPrintf_FormatsText(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Printf("hello %s, count=%d", "world", 42)
	output := buf.String()
	expected := fmt.Sprintf("hello %s, count=%d", "world", 42)
	assert.Equal(t, expected, output)
}

func TestAllMethods_NonTTY_NoAnsiEscapes(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Success("ok")
	p.Error("err")
	p.Warning("warn")
	p.Info("info")
	p.Header("head")
	p.Muted("muted")
	output := buf.String()
	// In non-TTY mode, there should be no ANSI escape sequences
	assert.NotContains(t, output, "\x1b[", "non-TTY output should not contain ANSI escape codes")
}
