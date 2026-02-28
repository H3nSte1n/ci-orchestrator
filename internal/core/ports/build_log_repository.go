package ports

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
)

type BuildLogRepository interface {
	Save(ctx context.Context, buildLog *domain.BuildLog) error
}
