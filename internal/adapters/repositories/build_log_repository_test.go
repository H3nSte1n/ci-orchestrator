package repositories

import (
	"context"
	"errors"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func buildLogTestData() domain.BuildLog {
	return domain.BuildLog{
		ID:      "ci-id build log",
		BuildID: "1",
		Stream:  domain.LogStderr,
		Content: "The first line",
	}
}

func TestBuildLogRepository_Save_Success(t *testing.T) {
	mockDB := new(mockDB)
	mockDB.Error = nil

	mockDB.On("WithContext", mock.Anything).Return(mockDB)
	mockDB.On("Create", mock.Anything).Return(mockDB)

	ctx := context.Background()
	buildLog := buildLogTestData()
	repo := buildLogRepository{db: mockDB}
	err := repo.Save(ctx, &buildLog)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestBuildLogRepository_Save_Error(t *testing.T) {
	mockDB := new(mockDB)
	expectedErr := errors.New("create failed")
	mockDB.Error = expectedErr

	mockDB.On("WithContext", mock.Anything).Return(mockDB)
	mockDB.On("Create", mock.Anything).Return(mockDB)

	ctx := context.Background()
	buildLog := buildLogTestData()
	repo := buildLogRepository{db: mockDB}
	err := repo.Save(ctx, &buildLog)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockDB.AssertExpectations(t)
}
