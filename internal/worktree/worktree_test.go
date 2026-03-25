package worktree_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xenciscbc/mysd/internal/worktree"
)

// initGitRepo creates a temp dir with a git repo that has one commit.
// This is required for git worktree add to work.
func initGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v failed: %s", args, out)
	}

	run("init")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "Test")
	// Create initial commit so worktree add works
	run("commit", "--allow-empty", "-m", "init")

	return dir
}

func TestToSlug(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"Setup Auth System", "setup-auth-system"},
		{"Setup Auth", "setup-auth"},
		{"Hello World!", "hello-world"},
		{"  leading trailing  ", "leading-trailing"},
		{"multiple---hyphens", "multiple-hyphens"},
		{"special!@#chars", "special-chars"},
		{"MixedCASE", "mixedcase"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := worktree.ToSlug(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestCreate_BranchName(t *testing.T) {
	repoRoot := initGitRepo(t)
	mgr := worktree.WorktreeManager{
		RepoRoot:    repoRoot,
		WorktreeDir: ".worktrees",
		ChangeName:  "mychange",
	}

	_, branch, err := mgr.Create(3, "Setup Auth")
	require.NoError(t, err)
	assert.Equal(t, "mysd/mychange/T3-setup-auth", branch)

	// Cleanup
	_ = mgr.Remove(3, branch)
}

func TestCreate_Path(t *testing.T) {
	repoRoot := initGitRepo(t)
	mgr := worktree.WorktreeManager{
		RepoRoot:    repoRoot,
		WorktreeDir: ".worktrees",
		ChangeName:  "mychange",
	}

	path, branch, err := mgr.Create(3, "Setup Auth")
	require.NoError(t, err)

	expected := filepath.Join(repoRoot, ".worktrees", "T3")
	assert.Equal(t, expected, path)

	// Cleanup
	_ = mgr.Remove(3, branch)
}

func TestCreate_ActualWorktree(t *testing.T) {
	repoRoot := initGitRepo(t)
	mgr := worktree.WorktreeManager{
		RepoRoot:    repoRoot,
		WorktreeDir: ".worktrees",
		ChangeName:  "mychange",
	}

	path, branch, err := mgr.Create(5, "Do Something")
	require.NoError(t, err)

	// Verify the worktree directory was actually created
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Cleanup
	_ = mgr.Remove(5, branch)
}

func TestRemove_CleansUp(t *testing.T) {
	repoRoot := initGitRepo(t)
	mgr := worktree.WorktreeManager{
		RepoRoot:    repoRoot,
		WorktreeDir: ".worktrees",
		ChangeName:  "mychange",
	}

	path, branch, err := mgr.Create(7, "Test Task")
	require.NoError(t, err)

	// Verify created
	_, err = os.Stat(path)
	require.NoError(t, err)

	// Remove
	err = mgr.Remove(7, branch)
	require.NoError(t, err)

	// Verify removed
	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err), "worktree dir should no longer exist")
}

func TestRemove_BranchDeleted(t *testing.T) {
	repoRoot := initGitRepo(t)
	mgr := worktree.WorktreeManager{
		RepoRoot:    repoRoot,
		WorktreeDir: ".worktrees",
		ChangeName:  "mychange",
	}

	_, branch, err := mgr.Create(8, "Branch Test")
	require.NoError(t, err)

	err = mgr.Remove(8, branch)
	require.NoError(t, err)

	// Check branch is gone
	cmd := exec.Command("git", "branch", "--list", branch)
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	require.NoError(t, err)
	assert.Empty(t, strings.TrimSpace(string(out)), "branch should be deleted after remove")
}

func TestCheckDiskSpace_Sufficient(t *testing.T) {
	repoRoot := initGitRepo(t)
	mgr := worktree.WorktreeManager{
		RepoRoot:    repoRoot,
		WorktreeDir: ".worktrees",
		ChangeName:  "test",
	}

	// 1 byte minimum — should always pass on any real system
	err := mgr.CheckDiskSpace(1)
	assert.NoError(t, err)
}

func TestCheckDiskSpace_Insufficient(t *testing.T) {
	repoRoot := initGitRepo(t)
	mgr := worktree.WorktreeManager{
		RepoRoot:    repoRoot,
		WorktreeDir: ".worktrees",
		ChangeName:  "test",
	}

	// Require more space than any system has (math.MaxUint64)
	const impossibleBytes = ^uint64(0) // max uint64
	err := mgr.CheckDiskSpace(impossibleBytes)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient disk space")
}

func TestCreate_WindowsLongPaths(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-only test")
	}

	repoRoot := initGitRepo(t)
	mgr := worktree.WorktreeManager{
		RepoRoot:    repoRoot,
		WorktreeDir: ".worktrees",
		ChangeName:  "mychange",
	}

	_, branch, err := mgr.Create(1, "Windows Test")
	require.NoError(t, err)

	// Check that core.longpaths was set
	cmd := exec.Command("git", "config", "core.longpaths")
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	require.NoError(t, err)
	assert.Equal(t, "true", strings.TrimSpace(string(out)))

	// Cleanup
	_ = mgr.Remove(1, branch)
}
