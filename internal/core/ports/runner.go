package ports

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
)

type Runner interface {
	Start(ctx context.Context, workdir string, command string, env []string) (<-chan domain.LogEvent, func() (int, error), error)
}
