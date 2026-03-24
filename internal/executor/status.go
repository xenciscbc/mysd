package executor

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
)

// StatusSummary aggregates workflow state, task progress, and requirement counts
// for display in the status dashboard.
type StatusSummary struct {
	ChangeName  string
	Phase       string
	TasksDone   int
	TasksTotal  int
	MustDone    int
	MustTotal   int
	ShouldDone  int
	ShouldTotal int
	MayTotal    int
	LastRun     string // formatted time, e.g. "2006-01-02 15:04" or "never"
}

// Lipgloss styles for the status dashboard.
var (
	styleHeader  = lipgloss.NewStyle().Bold(true)
	styleLabel   = lipgloss.NewStyle().Bold(true)
	styleDone    = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")) // green
	stylePending = lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b")) // amber
	styleMuted   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280")) // gray
	stylePhase   = lipgloss.NewStyle().Foreground(lipgloss.Color("#3b82f6")) // blue
)

// BuildStatusSummary computes a StatusSummary from the provided WorkflowState,
// task list, and requirement list.
func BuildStatusSummary(ws state.WorkflowState, tasks []spec.Task, reqs []spec.Requirement) StatusSummary {
	s := StatusSummary{
		ChangeName: ws.ChangeName,
		Phase:      string(ws.Phase),
	}

	// Compute task counts.
	for _, t := range tasks {
		s.TasksTotal++
		if t.Status == spec.StatusDone {
			s.TasksDone++
		}
	}

	// Compute requirement counts by keyword.
	for _, r := range reqs {
		switch r.Keyword {
		case spec.Must:
			s.MustTotal++
			if r.Status == spec.StatusDone {
				s.MustDone++
			}
		case spec.Should:
			s.ShouldTotal++
			if r.Status == spec.StatusDone {
				s.ShouldDone++
			}
		case spec.May:
			s.MayTotal++
		}
	}

	// Format LastRun.
	if ws.LastRun.IsZero() {
		s.LastRun = "never"
	} else {
		s.LastRun = ws.LastRun.Format("2006-01-02 15:04")
	}

	return s
}

// RenderStatus writes a lipgloss-styled status dashboard to w.
// The dashboard shows change name, workflow phase, task progress,
// MUST/SHOULD/MAY requirement counts, and last run time.
func RenderStatus(w io.Writer, summary StatusSummary) {
	// Header.
	header := fmt.Sprintf("=== %s ===", summary.ChangeName)
	fmt.Fprintln(w, styleHeader.Render(header))
	fmt.Fprintln(w)

	// Phase row.
	fmt.Fprintf(w, "%s  %s\n",
		styleLabel.Render("Phase:     "),
		stylePhase.Render(summary.Phase),
	)

	// Progress row with progress bar.
	progressBar := buildProgressBar(summary.TasksDone, summary.TasksTotal, 20)
	fmt.Fprintf(w, "%s  %s  %s\n",
		styleLabel.Render("Progress:  "),
		styleDone.Render(fmt.Sprintf("%d/%d tasks done", summary.TasksDone, summary.TasksTotal)),
		styleMuted.Render(progressBar),
	)

	// MUST row.
	mustPending := summary.MustTotal - summary.MustDone
	fmt.Fprintf(w, "%s  %s  %s  %s\n",
		styleLabel.Render("MUST items:"),
		styleDone.Render(fmt.Sprintf("%d done", summary.MustDone)),
		stylePending.Render(fmt.Sprintf("%d pending", mustPending)),
		styleMuted.Render(fmt.Sprintf("(%d total)", summary.MustTotal)),
	)

	// SHOULD row.
	shouldPending := summary.ShouldTotal - summary.ShouldDone
	fmt.Fprintf(w, "%s  %s  %s  %s\n",
		styleLabel.Render("SHOULD:    "),
		styleDone.Render(fmt.Sprintf("%d done", summary.ShouldDone)),
		stylePending.Render(fmt.Sprintf("%d pending", shouldPending)),
		styleMuted.Render(fmt.Sprintf("(%d total)", summary.ShouldTotal)),
	)

	// MAY row.
	fmt.Fprintf(w, "%s  %s\n",
		styleLabel.Render("MAY:       "),
		styleMuted.Render(fmt.Sprintf("%d noted", summary.MayTotal)),
	)

	// Last run row.
	fmt.Fprintf(w, "%s  %s\n",
		styleLabel.Render("Last run:  "),
		styleMuted.Render(summary.LastRun),
	)
}

// buildProgressBar creates a simple block-character progress bar.
// width is the total number of characters in the bar.
func buildProgressBar(done, total, width int) string {
	if total == 0 {
		return "[" + strings.Repeat("░", width) + "]"
	}
	filled := (done * width) / total
	if filled > width {
		filled = width
	}
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", width-filled) + "]"
}
