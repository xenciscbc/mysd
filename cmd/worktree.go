package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xenciscbc/mysd/internal/config"
	"github.com/xenciscbc/mysd/internal/spec"
	"github.com/xenciscbc/mysd/internal/state"
	"github.com/xenciscbc/mysd/internal/worktree"
)

var worktreeCmd = &cobra.Command{
	Use:   "worktree",
	Short: "Manage git worktrees for parallel execution",
	Long:  `Create and remove git worktrees for wave-parallel task execution.`,
}

var worktreeCreateCmd = &cobra.Command{
	Use:   "create <task_id> <task_name>",
	Short: "Create a worktree for a task",
	Args:  cobra.ExactArgs(2),
	RunE:  runWorktreeCreate,
}

var worktreeRemoveCmd = &cobra.Command{
	Use:   "remove <task_id> <branch>",
	Short: "Remove a worktree after merge",
	Args:  cobra.ExactArgs(2),
	RunE:  runWorktreeRemove,
}

func init() {
	worktreeCmd.AddCommand(worktreeCreateCmd)
	worktreeCmd.AddCommand(worktreeRemoveCmd)
	rootCmd.AddCommand(worktreeCmd)
}

// repoRoot finds the git repository root by running `git rev-parse --show-toplevel`.
func repoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("find repo root: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func runWorktreeCreate(cmd *cobra.Command, args []string) error {
	taskID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task_id %q: must be an integer", args[0])
	}
	taskName := args[1]

	cfg, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	root, err := repoRoot()
	if err != nil {
		return err
	}

	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return fmt.Errorf("detect spec dir: %w", err)
	}

	ws, err := state.LoadState(specDir)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	mgr := worktree.WorktreeManager{
		RepoRoot:    root,
		WorktreeDir: cfg.WorktreeDir,
		ChangeName:  ws.ChangeName,
	}

	path, branch, err := mgr.Create(taskID, taskName)
	if err != nil {
		return fmt.Errorf("create worktree: %w", err)
	}

	result := map[string]string{
		"path":   path,
		"branch": branch,
	}
	data, _ := json.Marshal(result)
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}

func runWorktreeRemove(cmd *cobra.Command, args []string) error {
	taskID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task_id %q: must be an integer", args[0])
	}
	branch := args[1]

	cfg, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	root, err := repoRoot()
	if err != nil {
		return err
	}

	mgr := worktree.WorktreeManager{
		RepoRoot:    root,
		WorktreeDir: cfg.WorktreeDir,
		ChangeName:  "", // Not needed for remove
	}

	if err := mgr.Remove(taskID, branch); err != nil {
		return fmt.Errorf("remove worktree: %w", err)
	}

	result := map[string]bool{"removed": true}
	data, _ := json.Marshal(result)
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}
