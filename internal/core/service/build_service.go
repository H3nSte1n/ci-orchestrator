package service

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"time"
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
	if err := s.buildRepo.Save(ctx, build); err != nil {
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

func (s *buildService) GetBuild(ctx context.Context, buildId string) (*domain.Build, error) {
	build, err := s.buildRepo.FindByID(ctx, buildId)
	if err != nil {
		return nil, err
	}

	return build, nil
}

func (s *buildService) ClaimNext(ctx context.Context, workerId string) (*domain.Build, error) {
	build, err := s.buildRepo.ClaimNext(ctx, workerId)

	if err != nil {
		return nil, err
	}

	return build, nil
}

func (s *buildService) CompleteBuild(ctx context.Context, buildId string, exitCode int, finishedAt *time.Time, error error) error {
	build := &domain.Build{
		ID:         buildId,
		Status:     domain.BuildStatusSuccess,
		ExitCode:   exitCode,
		FinishedAt: finishedAt,
	}

	if error != nil {
		build.Error = error.Error()
	}

	err := s.buildRepo.Update(ctx, build)

	return err
}
