package ports

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
)

type BuildRepository interface {
	Save(ctx context.Context, build *domain.Build) error
	Update(ctx context.Context, build *domain.Build) error
	FindByID(ctx context.Context, buildId string) (*domain.Build, error)
}
