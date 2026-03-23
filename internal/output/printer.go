package output

import (
	"fmt"
	"io"

	"github.com/charmbracelet/x/term"
)

// Printer provides styled terminal output with automatic TTY detection.
// In TTY mode, output is styled using lipgloss colors. In non-TTY mode
// (e.g., pipes, CI), plain text with prefixes is used instead.
type Printer struct {
	w     io.Writer
	isTTY bool
}

// NewPrinter creates a Printer for the given writer.
// TTY detection is performed by checking if the writer implements Fd() uintptr
// and calling term.IsTerminal.
func NewPrinter(w io.Writer) *Printer {
	isTTY := false
	type fder interface {
		Fd() uintptr
	}
	if f, ok := w.(fder); ok {
		isTTY = term.IsTerminal(f.Fd())
	}
	return &Printer{w: w, isTTY: isTTY}
}

// Success prints a success message. TTY: styled green bold. Non-TTY: "[OK] msg".
func (p *Printer) Success(msg string) {
	if p.isTTY {
		fmt.Fprintln(p.w, StyleSuccess.Render(msg))
	} else {
		fmt.Fprintf(p.w, "[OK] %s\n", msg)
	}
}

// Error prints an error message. TTY: styled red bold. Non-TTY: "[ERROR] msg".
func (p *Printer) Error(msg string) {
	if p.isTTY {
		fmt.Fprintln(p.w, StyleError.Render(msg))
	} else {
		fmt.Fprintf(p.w, "[ERROR] %s\n", msg)
	}
}

// Warning prints a warning message. TTY: styled amber. Non-TTY: "[WARN] msg".
func (p *Printer) Warning(msg string) {
	if p.isTTY {
		fmt.Fprintln(p.w, StyleWarning.Render(msg))
	} else {
		fmt.Fprintf(p.w, "[WARN] %s\n", msg)
	}
}

// Info prints an informational message. TTY: styled blue. Non-TTY: "[INFO] msg".
func (p *Printer) Info(msg string) {
	if p.isTTY {
		fmt.Fprintln(p.w, StyleInfo.Render(msg))
	} else {
		fmt.Fprintf(p.w, "[INFO] %s\n", msg)
	}
}

// Header prints a header message. TTY: styled bold underline. Non-TTY: "=== msg ===".
func (p *Printer) Header(msg string) {
	if p.isTTY {
		fmt.Fprintln(p.w, StyleHeader.Render(msg))
	} else {
		fmt.Fprintf(p.w, "=== %s ===\n", msg)
	}
}

// Muted prints a muted/dimmed message. TTY: styled gray. Non-TTY: plain msg.
func (p *Printer) Muted(msg string) {
	if p.isTTY {
		fmt.Fprintln(p.w, StyleMuted.Render(msg))
	} else {
		fmt.Fprintf(p.w, "%s\n", msg)
	}
}

// Printf writes a formatted string directly to the writer (no styling applied).
func (p *Printer) Printf(format string, args ...interface{}) {
	fmt.Fprintf(p.w, format, args...)
}
