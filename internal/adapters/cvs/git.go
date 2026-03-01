package vcs

import (
	"context"
	"fmt"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"os"
	"os/exec"
	"syscall"
)

type GitVCS struct{}

func NewGitVCS() ports.VCS {
	return &GitVCS{}
}

func (g *GitVCS) CloneAndCheckout(ctx context.Context, repoUrl, ref, destDir string) error {
	if err := g.runCMD(ctx, "", nil, "git", "clone", "--quiet", repoUrl, destDir); err != nil {
		return fmt.Errorf("git clone: %w", err)
	}

	if err := g.runCMD(ctx, destDir, nil, "git", "checkout", "--quiet", ref); err != nil {
		return fmt.Errorf("git checkout %q: %w", ref, err)
	}

	return nil
}

func (g *GitVCS) runCMD(ctx context.Context, dir string, env []string, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Env = append(os.Environ(), env...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return cmd.Run()
}
