package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
)

// WorktreeManager manages git worktrees for parallel task execution.
// Each task gets its own worktree at WorktreeDir/T{id}/ on branch mysd/{change}/T{id}-{slug}.
type WorktreeManager struct {
	RepoRoot    string // absolute path to repo root
	WorktreeDir string // relative directory name, e.g. ".worktrees"
	ChangeName  string // current change name for branch prefix
}

// Create adds a new worktree at .worktrees/T{id}/ on branch mysd/{change}/T{id}-{slug}.
// On Windows, also sets git config core.longpaths=true to handle long paths.
func (m *WorktreeManager) Create(id int, taskName string) (path string, branch string, err error) {
	// 1. Set longpaths on Windows (must be before any path operations)
	if runtime.GOOS == "windows" {
		if err := m.setLongPaths(); err != nil {
			return "", "", fmt.Errorf("set longpaths: %w", err)
		}
	}

	// 2. Check disk space (500MB threshold)
	if err := m.CheckDiskSpace(500 * 1024 * 1024); err != nil {
		return "", "", err
	}

	// 3. Compute paths and branch name
	slug := ToSlug(taskName)
	worktreePath := filepath.Join(m.RepoRoot, m.WorktreeDir, fmt.Sprintf("T%d", id))
	branchName := fmt.Sprintf("mysd/%s/T%d-%s", m.ChangeName, id, slug)

	// 4. git worktree add -b {branch} {path}
	cmd := exec.Command("git", "worktree", "add", "-b", branchName, worktreePath)
	cmd.Dir = m.RepoRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("git worktree add: %w\n%s", err, out)
	}

	return worktreePath, branchName, nil
}

// Remove deletes the worktree directory and branch (success cleanup).
// Uses --force to handle untracked files in the worktree.
// Branch delete failure is non-fatal (logs warning to stderr).
func (m *WorktreeManager) Remove(id int, branch string) error {
	worktreePath := filepath.Join(m.RepoRoot, m.WorktreeDir, fmt.Sprintf("T%d", id))

	// git worktree remove --force {path}
	rmCmd := exec.Command("git", "worktree", "remove", "--force", worktreePath)
	rmCmd.Dir = m.RepoRoot
	if out, err := rmCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git worktree remove: %w\n%s", err, out)
	}

	// git branch -d {branch} — non-fatal
	delCmd := exec.Command("git", "branch", "-d", branch)
	delCmd.Dir = m.RepoRoot
	if out, err := delCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: git branch -d %s: %v\n%s\n", branch, err, out)
	}

	// git worktree prune — cleanup stale metadata
	pruneCmd := exec.Command("git", "worktree", "prune")
	pruneCmd.Dir = m.RepoRoot
	if out, err := pruneCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: git worktree prune: %v\n%s\n", err, out)
	}

	return nil
}

// CheckDiskSpace returns an error if available bytes at RepoRoot < minBytes.
// If available space cannot be determined, returns nil (non-fatal).
func (m *WorktreeManager) CheckDiskSpace(minBytes uint64) error {
	available, err := getAvailableBytes(m.RepoRoot)
	if err != nil {
		// Non-fatal: can't determine disk space
		return nil
	}
	if available < minBytes {
		return fmt.Errorf(
			"insufficient disk space: need at least %dMB, have %dMB available at %s",
			minBytes/(1024*1024), available/(1024*1024), m.RepoRoot,
		)
	}
	return nil
}

// setLongPaths sets git config core.longpaths true for the repo.
// Required on Windows to handle paths longer than MAX_PATH.
func (m *WorktreeManager) setLongPaths() error {
	cmd := exec.Command("git", "config", "core.longpaths", "true")
	cmd.Dir = m.RepoRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git config core.longpaths true: %w\n%s", err, out)
	}
	return nil
}

// ToSlug converts a task name to a URL-safe slug suitable for branch names.
// "Setup Auth System" -> "setup-auth-system"
// Exported so tests can verify slug logic independently.
func ToSlug(name string) string {
	// Lowercase first
	name = strings.ToLower(name)

	// Replace non-alphanumeric with hyphen
	var b strings.Builder
	prevHyphen := false
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prevHyphen = false
		} else {
			if !prevHyphen && b.Len() > 0 {
				b.WriteByte('-')
				prevHyphen = true
			}
		}
	}

	// Trim trailing hyphen
	result := strings.TrimRight(b.String(), "-")
	return result
}
