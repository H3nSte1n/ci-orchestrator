package ports

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"time"
)

type BuildService interface {
	CreateBuild(ctx context.Context, build *domain.Build) error
	CancelBuild(ctx context.Context, buildId string) error
	UpdateStatus(ctx context.Context, buildId string, status domain.BuildStatus) error
	GetBuild(ctx context.Context, buildId string) (*domain.Build, error)
	ClaimNext(ctx context.Context, workerId string) (*domain.Build, error)
	CompleteBuild(ctx context.Context, buildId string, exitCode int, finishedAt *time.Time, error error) error
}
