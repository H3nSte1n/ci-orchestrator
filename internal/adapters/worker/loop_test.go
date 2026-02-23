package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestData() *domain.Build {
	return &domain.Build{
		ID:      "ci-id",
		RepoUrl: "https://github.com/test/repo",
		Ref:     "main",
		Command: "npm test",
		Status:  domain.BuildStatusPending,
	}
}

type mockBuildService struct {
	mock.Mock
	Error error
}

func (m *mockBuildService) ClaimNext(ctx context.Context, workerId string) (*domain.Build, error) {
	args := m.Called(ctx, workerId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Build), args.Error(1)
}

func (m *mockBuildService) CreateBuild(ctx context.Context, build *domain.Build) error {
	args := m.Called(ctx, build)
	return args.Error(0)
}

func (m *mockBuildService) CancelBuild(ctx context.Context, buildId string) error {
	args := m.Called(ctx, buildId)
	return args.Error(0)
}

func (m *mockBuildService) UpdateStatus(ctx context.Context, buildId string, status domain.BuildStatus) error {
	args := m.Called(ctx, buildId, status)
	return args.Error(0)
}

func (m *mockBuildService) GetBuild(ctx context.Context, buildId string) (*domain.Build, error) {
	args := m.Called(ctx, buildId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Build), args.Error(1)
}

func (m *mockBuildService) GetError() error {
	return m.Error
}

func TestWorker_ClaimAndProcess_Success(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("ClaimNext", mock.Anything, "worker-1").Return(buildTestData(), nil)

	worker := NewWorker("worker-1", mockBuildService, 100*time.Millisecond)
	err := worker.claimAndProcess(context.Background())

	assert.NoError(t, err)
	mockBuildService.AssertCalled(t, "ClaimNext", mock.Anything, "worker-1")
}

func TestWorker_ClaimAndProcess_Error(t *testing.T) {
	mockBuildService := new(mockBuildService)
	expectedErr := errors.New("db error")
	mockBuildService.On("ClaimNext", mock.Anything, "worker-1").Return(nil, expectedErr)

	worker := NewWorker("worker-1", mockBuildService, 100*time.Millisecond)
	err := worker.claimAndProcess(context.Background())

	assert.ErrorIs(t, err, expectedErr)
	mockBuildService.AssertCalled(t, "ClaimNext", mock.Anything, "worker-1")
}

func TestWorker_ClaimAndProcess_NoBuilds(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("ClaimNext", mock.Anything, "worker-1").Return(nil, nil)

	worker := NewWorker("worker-1", mockBuildService, 100*time.Millisecond)
	err := worker.claimAndProcess(context.Background())

	assert.NoError(t, err)
	mockBuildService.AssertCalled(t, "ClaimNext", mock.Anything, "worker-1")
}

func TestWorker_Run_ExitsOnContextCancel(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("ClaimNext", mock.Anything, mock.Anything).Return(buildTestData(), nil)

	worker := NewWorker("worker-1", mockBuildService, 100*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	err := worker.Run(ctx)
	assert.Error(t, err)
}
