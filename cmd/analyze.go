package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/xenciscbc/mysd/internal/analyzer"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
)

var analyzeJSON bool

var analyzeCmd = &cobra.Command{
	Use:   "analyze [change-name]",
	Short: "Cross-artifact structural analysis",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runAnalyze,
}

func init() {
	analyzeCmd.Flags().BoolVar(&analyzeJSON, "json", false, "output as JSON")
	rootCmd.AddCommand(analyzeCmd)
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return fmt.Errorf("no spec directory: %w", err)
	}

	var changeName string
	if len(args) > 0 {
		changeName = args[0]
	} else {
		ws, err := state.LoadState(specDir)
		if err != nil {
			return fmt.Errorf("no active change and no change name provided: %w", err)
		}
		if ws.ChangeName == "" {
			return fmt.Errorf("no active change: provide a change name or run 'mysd propose' first")
		}
		changeName = ws.ChangeName
	}

	changeDir := fmt.Sprintf("%s/changes/%s", specDir, changeName)
	result := analyzer.Analyze(changeDir)

	if analyzeJSON {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal result: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	printAnalyzeSummary(cmd.OutOrStdout(), result)
	return nil
}

func printAnalyzeSummary(out io.Writer, result analyzer.AnalysisResult) {
	criticalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	suggestionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	cleanStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)

	fmt.Fprintf(out, "\nAnalysis: %s\n", result.ChangeID)
	fmt.Fprintln(out, "──────────────────────────────")

	for _, dim := range result.Dimensions {
		style := cleanStyle
		if dim.FindingCount > 0 {
			style = warningStyle
		}
		fmt.Fprintf(out, "%-15s %s\n", dim.Dimension+":", style.Render(dim.Status))
	}

	if len(result.Findings) > 0 {
		fmt.Fprintln(out, "\nFindings:")
		for _, f := range result.Findings {
			var style lipgloss.Style
			switch f.Severity {
			case string(analyzer.SeverityCritical):
				style = criticalStyle
			case string(analyzer.SeverityWarning):
				style = warningStyle
			default:
				style = suggestionStyle
			}
			fmt.Fprintf(out, "  %s [%s] %s\n", style.Render(f.ID), f.Location, f.Summary)
		}
	}

	fmt.Fprintln(out)
}
