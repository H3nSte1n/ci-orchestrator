package service

import (
	"context"
	"errors"
	"testing"

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

type MockBuildRepository struct {
	mock.Mock
}

func (m *MockBuildRepository) Save(ctx context.Context, build *domain.Build) error {
	args := m.Called(ctx, build)
	return args.Error(0)
}

func (m *MockBuildRepository) Update(ctx context.Context, build *domain.Build) error {
	args := m.Called(ctx, build)
	return args.Error(0)
}

func (m *MockBuildRepository) FindByID(ctx context.Context, buildId string) (*domain.Build, error) {
	args := m.Called(ctx, buildId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Build), args.Error(1)
}

func (m *MockBuildRepository) ClaimNext(ctx context.Context, workerId string) (*domain.Build, error) {
	args := m.Called(ctx, workerId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Build), args.Error(1)
}

func TestBuildService_CreateBuild_Success(t *testing.T) {
	mockRepo := new(MockBuildRepository)
	ctx := context.Background()
	build := buildTestData()
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	service := NewBuildService(mockRepo)
	err := service.CreateBuild(ctx, build)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBuildService_CreateBuild_Error(t *testing.T) {
	mockRepo := new(MockBuildRepository)
	ctx := context.Background()
	build := buildTestData()

	expectedErr := errors.New("database error")
	mockRepo.On("Save", mock.Anything, mock.Anything).Return(expectedErr)

	service := NewBuildService(mockRepo)
	err := service.CreateBuild(ctx, build)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestBuildService_CancelBuild_Success(t *testing.T) {
	mockRepo := new(MockBuildRepository)
	ctx := context.Background()
	buildId := "test-build-id"
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	service := NewBuildService(mockRepo)
	err := service.CancelBuild(ctx, buildId)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBuildService_CancelBuild_Error(t *testing.T) {
	mockRepo := new(MockBuildRepository)
	ctx := context.Background()
	buildId := "test-build-id"

	expectedErr := errors.New("failed to update build")
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(expectedErr)

	service := NewBuildService(mockRepo)
	err := service.CancelBuild(ctx, buildId)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestBuildService_UpdateStatus_Success(t *testing.T) {
	mockRepo := new(MockBuildRepository)
	ctx := context.Background()
	buildId := "test-build-id"
	newStatus := domain.BuildStatusRunning

	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	service := NewBuildService(mockRepo)
	err := service.UpdateStatus(ctx, buildId, newStatus)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBuildService_UpdateStatus_Error(t *testing.T) {
	mockRepo := new(MockBuildRepository)
	ctx := context.Background()
	buildId := "test-build-id"
	newStatus := domain.BuildStatusFailed

	expectedErr := errors.New("database error")
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(expectedErr)

	service := NewBuildService(mockRepo)
	err := service.UpdateStatus(ctx, buildId, newStatus)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestBuildService_GetBuild_Success(t *testing.T) {
	mockRepo := new(MockBuildRepository)
	ctx := context.Background()
	buildId := "test-build-id"
	expectedBuild := buildTestData()

	mockRepo.On("FindByID", mock.Anything, mock.Anything).Return(expectedBuild, nil)

	service := NewBuildService(mockRepo)
	build, err := service.GetBuild(ctx, buildId)

	assert.NoError(t, err)
	assert.Equal(t, expectedBuild, build)
}

func TestBuildService_GetBuild_Error(t *testing.T) {
	mockRepo := new(MockBuildRepository)
	ctx := context.Background()
	buildId := "test-build-id"
	expectedErr := errors.New("build not found")

	mockRepo.On("FindByID", mock.Anything, mock.Anything).Return(nil, expectedErr)

	service := NewBuildService(mockRepo)
	_, err := service.GetBuild(ctx, buildId)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
