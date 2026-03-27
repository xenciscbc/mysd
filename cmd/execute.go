package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/xenciscbc/mysd/internal/config"
	"github.com/xenciscbc/mysd/internal/executor"
	"github.com/xenciscbc/mysd/internal/output"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/spf13/cobra"
)

var contextOnly bool

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Run tasks with pre-execution alignment",
	RunE:  runExecute,
}

func init() {
	executeCmd.Flags().BoolVar(&contextOnly, "context-only", false, "output execution context as JSON (for SKILL.md consumption)")
	rootCmd.AddCommand(executeCmd)
}

func runExecute(cmd *cobra.Command, args []string) error {
	p := output.NewPrinter(cmd.OutOrStdout())

	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return err
	}

	ws, err := state.LoadState(specDir)
	if err != nil {
		return err
	}

	cfg, err := config.Load(".")
	if err != nil {
		return err
	}

	// Override config with flags
	if cmd.Flags().Changed("tdd") {
		cfg.TDD, _ = cmd.Flags().GetBool("tdd")
	}
	if cmd.Flags().Changed("atomic-commits") {
		cfg.AtomicCommits, _ = cmd.Flags().GetBool("atomic-commits")
	}
	if cmd.Flags().Changed("execution-mode") {
		cfg.ExecutionMode, _ = cmd.Flags().GetString("execution-mode")
	}
	if cmd.Flags().Changed("agent-count") {
		cfg.AgentCount, _ = cmd.Flags().GetInt("agent-count")
	}

	if contextOnly {
		ctx, err := executor.BuildContext(specDir, ws.ChangeName, cfg)
		if err != nil {
			return err
		}
		data, _ := json.MarshalIndent(ctx, "", "  ")
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	// Non-context-only: print guidance directing to SKILL.md
	p.Info("Use /mysd:apply in Claude Code to run with AI alignment gate")
	p.Info("Or use: mysd execute --context-only | jq")
	return nil
}
