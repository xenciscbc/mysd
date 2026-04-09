package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/xenciscbc/mysd/internal/validator"
)

var validateJSON bool

var validateCmd = &cobra.Command{
	Use:   "validate [change-name]",
	Short: "Validate artifact structure and required fields",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runValidate,
}

func init() {
	validateCmd.Flags().BoolVar(&validateJSON, "json", false, "output as JSON")
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
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
	result := validator.Validate(changeDir)

	if validateJSON {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal result: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	printValidateSummary(cmd.OutOrStdout(), result)
	return nil
}

func printValidateSummary(out io.Writer, result validator.ValidationResult) {
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	passStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)

	fmt.Fprintf(out, "\nValidate: %s\n", result.ChangeID)
	fmt.Fprintln(out, "──────────────────────────────")

	for _, e := range result.Errors {
		fmt.Fprintf(out, "  %s %s: %s\n", errorStyle.Render("✗"), e.Location, e.Message)
	}
	for _, w := range result.Warnings {
		fmt.Fprintf(out, "  %s %s: %s\n", warnStyle.Render("⚠"), w.Location, w.Message)
	}

	errorCount := len(result.Errors)
	warnCount := len(result.Warnings)

	if errorCount == 0 && warnCount == 0 {
		fmt.Fprintf(out, "\n%s\n\n", passStyle.Render("Result: PASS"))
	} else if errorCount == 0 {
		fmt.Fprintf(out, "\n%s (%d warning(s))\n\n", passStyle.Render("Result: PASS"), warnCount)
	} else {
		fmt.Fprintf(out, "\n%s (%d error(s), %d warning(s))\n\n", errorStyle.Render("Result: FAIL"), errorCount, warnCount)
	}
}
