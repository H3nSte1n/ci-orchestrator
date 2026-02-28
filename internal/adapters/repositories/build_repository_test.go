package repositories

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

func TestBuildRepository_Save_Success(t *testing.T) {
	mockDB := new(mockDB)
	mockDB.Error = nil

	mockDB.On("WithContext", mock.Anything).Return(mockDB)
	mockDB.On("Create", mock.Anything).Return(mockDB)

	repo := &buildRepository{db: mockDB}
	ctx := context.Background()
	build := buildTestData()
	err := repo.Save(ctx, build)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestBuildRepository_Save_Error(t *testing.T) {
	mockDB := new(mockDB)
	expectedErr := errors.New("create failed")
	mockDB.Error = expectedErr

	mockDB.On("WithContext", mock.Anything).Return(mockDB)
	mockDB.On("Create", mock.Anything).Return(mockDB)

	repo := &buildRepository{db: mockDB}
	ctx := context.Background()
	build := buildTestData()
	err := repo.Save(ctx, build)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockDB.AssertExpectations(t)
}

func TestBuildRepository_Update_Success(t *testing.T) {
	mockDB := new(mockDB)
	mockDB.Error = nil

	mockDB.On("WithContext", mock.Anything).Return(mockDB)
	mockDB.On("Where", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Updates", mock.Anything).Return(mockDB)

	repo := &buildRepository{db: mockDB}
	ctx := context.Background()
	build := buildTestData()
	err := repo.Update(ctx, build)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestBuildRepository_Update_Error(t *testing.T) {
	mockDB := new(mockDB)
	expectedErr := errors.New("update failed")
	mockDB.Error = expectedErr

	mockDB.On("WithContext", mock.Anything).Return(mockDB)
	mockDB.On("Where", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Updates", mock.Anything).Return(mockDB)

	repo := &buildRepository{db: mockDB}
	ctx := context.Background()
	build := buildTestData()
	err := repo.Update(ctx, build)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockDB.AssertExpectations(t)
}

func TestBuildRepository_FindByID_Success(t *testing.T) {
	mockDB := new(mockDB)
	mockDB.Error = nil

	mockDB.On("WithContext", mock.Anything).Return(mockDB)
	mockDB.On("Where", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("First", mock.Anything).Return(mockDB).Run(func(args mock.Arguments) {
		build := args.Get(0).(*domain.Build)
		*build = *buildTestData()
	})

	repo := &buildRepository{db: mockDB}
	ctx := context.Background()
	build, err := repo.FindByID(ctx, "ci-id")

	assert.NoError(t, err)
	assert.Equal(t, "ci-id", build.ID)
	mockDB.AssertExpectations(t)
}

func TestBuildRepository_FindByID_Error(t *testing.T) {
	mockDB := new(mockDB)
	expectedErr := errors.New("database error")
	mockDB.Error = expectedErr

	mockDB.On("WithContext", mock.Anything).Return(mockDB)
	mockDB.On("Where", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("First", mock.Anything).Return(mockDB)

	repo := &buildRepository{db: mockDB}
	ctx := context.Background()
	_, err := repo.FindByID(ctx, "ci-id")

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockDB.AssertExpectations(t)
}

func TestBuildRepository_ClaimNext_Success(t *testing.T) {
	mockDB := new(mockDB)
	mockDB.Error = nil

	mockDB.On("WithContext", mock.Anything).Return(mockDB)
	mockDB.On("Where", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Order", mock.Anything).Return(mockDB)
	mockDB.On("Clauses", mock.Anything).Return(mockDB)
	mockDB.On("First", mock.Anything).Return(mockDB).Run(func(args mock.Arguments) {
		build := args.Get(0).(*domain.Build)
		*build = *buildTestData()
	})
	mockDB.On("Model", mock.Anything).Return(mockDB)
	mockDB.On("Updates", mock.Anything).Return(mockDB)

	repo := &buildRepository{db: mockDB}
	ctx := context.Background()
	build, err := repo.ClaimNext(ctx, "worker-id")

	assert.NoError(t, err)
	assert.Equal(t, "ci-id", build.ID)
	mockDB.AssertExpectations(t)
}

func TestBuildRepository_ClaimNext_Error(t *testing.T) {
	mockDB := new(mockDB)
	expectedErr := errors.New("database error")
	mockDB.Error = expectedErr

	mockDB.On("WithContext", mock.Anything).Return(mockDB)
	mockDB.On("Where", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("Order", mock.Anything).Return(mockDB)
	mockDB.On("Clauses", mock.Anything).Return(mockDB)
	mockDB.On("First", mock.Anything).Return(mockDB)

	repo := &buildRepository{db: mockDB}
	ctx := context.Background()
	build, err := repo.ClaimNext(ctx, "worker-id")

	assert.Error(t, err)
	assert.Nil(t, build)
	assert.Equal(t, expectedErr, err)
	mockDB.AssertExpectations(t)
}
