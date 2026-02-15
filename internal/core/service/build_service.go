package service

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
)

type buildService struct {
	buildRepo ports.BuildRepository
}

func NewBuildService(buildRepository ports.BuildRepository) ports.BuildService {
	return &buildService{
		buildRepo: buildRepository,
	}
}

func (s *buildService) CreateBuild(ctx context.Context, build *domain.Build) error {
	err := s.buildRepo.Save(ctx, build)
	if err != nil {
		return err
	}

	return nil
}

func (s *buildService) CancelBuild(ctx context.Context, buildId string) error {
	err := s.buildRepo.Update(ctx, &domain.Build{
		ID:     buildId,
		Status: domain.BuildStatusCanceled,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *buildService) UpdateStatus(ctx context.Context, buildId string, status domain.BuildStatus) error {
	err := s.buildRepo.Update(ctx, &domain.Build{
		ID:     buildId,
		Status: status,
	})

	if err != nil {
		return err
	}

	return nil
}
