package ports

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
)

type BuildLogService interface {
	AppendLog(ctx context.Context, buildId string, logEvent domain.LogEvent) error
}
