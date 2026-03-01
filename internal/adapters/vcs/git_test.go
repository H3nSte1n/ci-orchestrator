package vcs

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func run(t *testing.T, dir string, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "cmd failed: %s %v\n%s", name, args, string(out))
	return string(out)
}

func createRepo(t *testing.T, src string) string {
	run(t, src, "git", "init")
	run(t, src, "git", "config", "user.email", "test@example.com")
	run(t, src, "git", "config", "user.name", "test")

	require.NoError(t, os.WriteFile(filepath.Join(src, "file.txt"), []byte("v1\n"), 0o644))
	run(t, src, "git", "add", ".")
	run(t, src, "git", "commit", "-m", "c1")

	return strings.TrimSpace(run(t, src, "git", "rev-parse", "HEAD"))
}

func TestGitVCS_CloneAndCheckout_LocalRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	tmp := t.TempDir()
	src := filepath.Join(tmp, "src")
	dest := filepath.Join(tmp, "dest")
	require.NoError(t, os.MkdirAll(src, 0o755))

	sha := createRepo(t, src)
	g := &GitVCS{}
	err := g.CloneAndCheckout(context.Background(), src, sha, dest)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(dest, "file.txt"))
	require.NoError(t, err)

	head := strings.TrimSpace(run(t, dest, "git", "rev-parse", "HEAD"))
	require.Equal(t, sha, head)
}

func TestNewBuildLogRepository(t *testing.T) {
	repo := NewGitVCS()

	assert.NotNil(t, repo)
	assert.Implements(t, (*ports.VCS)(nil), repo)
}
