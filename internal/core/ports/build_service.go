package ports

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
)

type BuildService interface {
	CreateBuild(ctx context.Context, build *domain.Build) error
	CancelBuild(ctx context.Context, buildId string) error
	UpdateStatus(ctx context.Context, buildId string, status domain.BuildStatus) error
	GetBuild(ctx context.Context, buildId string) (*domain.Build, error)
}
