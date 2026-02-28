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

type mockBuildLogService struct {
	mock.Mock
}

func (m *mockBuildLogService) AppendLog(ctx context.Context, buildId string, ev domain.LogEvent) error {
	args := m.Called(ctx, buildId, ev)
	return args.Error(0)
}

type stubRunner struct {
	exitCode int
	runErr   error
	events   []domain.LogEvent
}

func (r *stubRunner) Start(_ context.Context, _, _ string, _ []string) (<-chan domain.LogEvent, func() (int, error), error) {
	ch := make(chan domain.LogEvent, len(r.events))
	for _, e := range r.events {
		ch <- e
	}
	close(ch)

	waitFn := func() (int, error) {
		return r.exitCode, r.runErr
	}

	return ch, waitFn, nil
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

func (m *mockBuildService) CompleteBuild(ctx context.Context, buildId string, exitCode int, finishedAt *time.Time, error error) error {
	args := m.Called(ctx, buildId, exitCode, finishedAt, error)
	return args.Error(0)
}

func (m *mockBuildService) GetError() error {
	return m.Error
}

func TestWorker_ClaimAndProcess_Success(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("ClaimNext", mock.Anything, "worker-1").Return(buildTestData(), nil)
	mockBuildService.On("CompleteBuild", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockBuildLogService := new(mockBuildLogService)
	mockBuildLogService.On("AppendLog", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	runner := &stubRunner{exitCode: 0, runErr: nil, events: []domain.LogEvent{{Stream: domain.LogStdout, Line: "hello", Time: time.Now()}}}

	worker := NewWorker("worker-1", mockBuildService, mockBuildLogService, 100*time.Millisecond, runner)
	err := worker.claimAndProcess(context.Background())

	assert.NoError(t, err)
	mockBuildService.AssertCalled(t, "ClaimNext", mock.Anything, "worker-1")
	mockBuildLogService.AssertCalled(t, "AppendLog", mock.Anything, "ci-id", mock.MatchedBy(func(e domain.LogEvent) bool {
		return e.Stream == domain.LogStdout && e.Line == "hello"
	}))
	mockBuildLogService.AssertNumberOfCalls(t, "AppendLog", 1)
}

func TestWorker_ClaimAndProcess_Error(t *testing.T) {
	mockBuildService := new(mockBuildService)
	expectedErr := errors.New("db error")
	mockBuildService.On("ClaimNext", mock.Anything, "worker-1").Return(nil, expectedErr)

	mockBuildLogService := new(mockBuildLogService)
	runner := &stubRunner{exitCode: 0, runErr: nil}

	worker := NewWorker("worker-1", mockBuildService, mockBuildLogService, 100*time.Millisecond, runner)
	err := worker.claimAndProcess(context.Background())

	assert.ErrorIs(t, err, expectedErr)
	mockBuildService.AssertCalled(t, "ClaimNext", mock.Anything, "worker-1")
}

func TestWorker_ClaimAndProcess_NoBuilds(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("ClaimNext", mock.Anything, "worker-1").Return(nil, nil)

	mockBuildLogService := new(mockBuildLogService)
	runner := &stubRunner{exitCode: 0, runErr: nil}

	worker := NewWorker("worker-1", mockBuildService, mockBuildLogService, 100*time.Millisecond, runner)
	err := worker.claimAndProcess(context.Background())

	assert.NoError(t, err)
	mockBuildService.AssertCalled(t, "ClaimNext", mock.Anything, "worker-1")
}

func TestWorker_ClaimAndProcess_CompleteBuildError(t *testing.T) {
	mockBuildService := new(mockBuildService)
	expectedErr := errors.New("db error")
	mockBuildService.On("ClaimNext", mock.Anything, "worker-1").Return(buildTestData(), nil)
	mockBuildService.On("CompleteBuild", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedErr)

	mockBuildLogService := new(mockBuildLogService)
	runner := &stubRunner{exitCode: 0, runErr: nil}

	worker := NewWorker("worker-1", mockBuildService, mockBuildLogService, 100*time.Millisecond, runner)
	err := worker.claimAndProcess(context.Background())

	assert.ErrorIs(t, err, expectedErr)
	mockBuildService.AssertCalled(t, "ClaimNext", mock.Anything, "worker-1")
	mockBuildService.AssertCalled(t, "CompleteBuild", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestWorker_Run_ExitsOnContextCancel(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("ClaimNext", mock.Anything, mock.Anything).Return(nil, nil)
	mockBuildLogService := new(mockBuildLogService)
	runner := &stubRunner{exitCode: 0, runErr: nil}

	worker := NewWorker("worker-1", mockBuildService, mockBuildLogService, 100*time.Millisecond, runner)
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	err := worker.Run(ctx)
	assert.Error(t, err)
}
