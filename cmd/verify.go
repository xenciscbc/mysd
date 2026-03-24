package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/xenciscbc/mysd/internal/roadmap"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/xenciscbc/mysd/internal/verifier"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Goal-backward verification of MUST items",
	RunE:  runVerify,
}

func init() {
	verifyCmd.Flags().Bool("context-only", false, "Output verification context as JSON (for /mysd:verify agent consumption)")
	verifyCmd.Flags().String("write-results", "", "Path to verifier report JSON to process")
	rootCmd.AddCommand(verifyCmd)
}

func runVerify(cmd *cobra.Command, args []string) error {
	contextOnly, _ := cmd.Flags().GetBool("context-only")
	writeResults, _ := cmd.Flags().GetString("write-results")

	if !contextOnly && writeResults == "" {
		return runVerifyNoFlags(cmd)
	}

	specsDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return err
	}

	ws, err := state.LoadState(specsDir)
	if err != nil {
		return err
	}

	if contextOnly {
		return runVerifyContextOnly(cmd.OutOrStdout(), specsDir, ws)
	}

	// write-results path
	return runVerifyWriteResults(cmd.OutOrStdout(), specsDir, &ws, writeResults)
}

// runVerifyNoFlags returns a usage hint error when no flags are provided.
func runVerifyNoFlags(cmd *cobra.Command) error {
	return fmt.Errorf("usage: mysd verify --context-only OR mysd verify --write-results <path>\nDirect execution via /mysd:verify")
}

// runVerifyContextOnly builds and outputs the VerificationContext JSON to out.
// Returns an error if there is no active change (empty ChangeName).
func runVerifyContextOnly(out io.Writer, specsDir string, ws state.WorkflowState) error {
	if ws.ChangeName == "" {
		return fmt.Errorf("no active change: run 'mysd propose' to start a new change")
	}

	ctx, err := verifier.BuildVerificationContext(specsDir, ws.ChangeName)
	if err != nil {
		return fmt.Errorf("build verification context: %w", err)
	}

	data, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal verification context: %w", err)
	}

	_, err = out.Write(data)
	return err
}

// runVerifyWriteResults reads a verifier report JSON, writes verification artifacts,
// updates verification-status.json, and transitions state if MUST all pass.
func runVerifyWriteResults(out io.Writer, specsDir string, ws *state.WorkflowState, reportPath string) error {
	// 1. Read and parse report JSON
	data, err := os.ReadFile(reportPath)
	if err != nil {
		return fmt.Errorf("read report file: %w", err)
	}

	report, err := verifier.ParseVerifierReport(data)
	if err != nil {
		return fmt.Errorf("parse verifier report: %w", err)
	}

	changeDir := filepath.Join(specsDir, "changes", ws.ChangeName)

	// 2. Write verification.md
	if err := verifier.WriteVerificationReport(changeDir, report); err != nil {
		return fmt.Errorf("write verification report: %w", err)
	}

	// 3. Write gap-report.md if there are failures
	if err := verifier.WriteGapReport(changeDir, report); err != nil {
		return fmt.Errorf("write gap report: %w", err)
	}

	// 4. Update verification-status.json — MUST items only
	vs := spec.VerificationStatus{
		ChangeName:   ws.ChangeName,
		VerifiedAt:   time.Now().UTC(),
		Requirements: make(map[string]spec.ItemStatus),
	}
	for _, r := range report.Results {
		if r.Keyword == "MUST" {
			if r.Pass {
				vs.Requirements[r.ID] = spec.StatusDone
			} else {
				vs.Requirements[r.ID] = spec.StatusBlocked
			}
		}
	}
	if err := spec.WriteVerificationStatus(changeDir, vs); err != nil {
		return fmt.Errorf("write verification status: %w", err)
	}

	// 5. State transition
	pass := report.MustPass
	ws.VerifyPass = &pass
	if report.MustPass {
		if err := state.Transition(ws, state.PhaseVerified); err != nil {
			return fmt.Errorf("state transition: %w", err)
		}
	}
	if err := state.SaveState(specsDir, *ws); err != nil {
		return fmt.Errorf("save state: %w", err)
	}
	if trackErr := roadmap.UpdateTracking(specsDir, *ws); trackErr != nil {
		fmt.Fprintf(os.Stderr, "warning: roadmap tracking update failed: %v\n", trackErr)
	}

	// 6. Terminal summary — lipgloss styled
	printVerifySummary(out, report)

	return nil
}

// printVerifySummary prints a lipgloss-styled summary of verification results.
func printVerifySummary(out io.Writer, report verifier.VerifierReport) {
	passStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true) // green
	failStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)  // red
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))              // gray

	var mustTotal, mustPassed, shouldTotal, shouldPassed, mayTotal, mayPassed int
	for _, r := range report.Results {
		switch r.Keyword {
		case "MUST":
			mustTotal++
			if r.Pass {
				mustPassed++
			}
		case "SHOULD":
			shouldTotal++
			if r.Pass {
				shouldPassed++
			}
		case "MAY":
			mayTotal++
			if r.Pass {
				mayPassed++
			}
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Verification Results")
	fmt.Fprintln(out, "--------------------")

	// MUST line
	mustLabel := passStyle.Render("PASS")
	if mustTotal > 0 && mustPassed < mustTotal {
		mustLabel = failStyle.Render("FAIL")
	}
	fmt.Fprintf(out, "MUST:   [%s] %d/%d passed\n", mustLabel, mustPassed, mustTotal)

	// SHOULD line
	shouldLabel := passStyle.Render("PASS")
	if shouldTotal > 0 && shouldPassed < shouldTotal {
		shouldLabel = dimStyle.Render("WARN")
	}
	if shouldTotal == 0 {
		shouldLabel = dimStyle.Render("N/A ")
	}
	fmt.Fprintf(out, "SHOULD: [%s] %d/%d passed\n", shouldLabel, shouldPassed, shouldTotal)

	// MAY line
	mayLabel := dimStyle.Render("INFO")
	fmt.Fprintf(out, "MAY:    [%s] %d/%d passed\n", mayLabel, mayPassed, mayTotal)

	fmt.Fprintln(out)
	if report.MustPass {
		fmt.Fprintln(out, passStyle.Render("All MUST items satisfied — state transitioned to verified"))
	} else {
		fmt.Fprintln(out, failStyle.Render("MUST failures detected — see gap-report.md for details"))
	}
}
