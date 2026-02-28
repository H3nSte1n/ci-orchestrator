package service

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
)

type buildLogService struct {
	buildLogRepo ports.BuildLogRepository
}

func NewBuildLogService(repo ports.BuildLogRepository) ports.BuildLogService {
	return &buildLogService{
		buildLogRepo: repo,
	}
}

func (s *buildLogService) AppendLog(ctx context.Context, buildId string, logEvent domain.LogEvent) error {
	b := domain.BuildLog{
		BuildID: buildId,
		Stream:  logEvent.Stream,
		Content: logEvent.Line,
	}

	if err := s.buildLogRepo.Save(ctx, &b); err != nil {
		return err
	}

	return nil
}
