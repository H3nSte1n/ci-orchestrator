package service

import (
	"context"
	"errors"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockBuildLogRepository struct {
	mock.Mock
}

func (m *mockBuildLogRepository) Save(ctx context.Context, buildLog *domain.BuildLog) error {
	args := m.Called(ctx, buildLog)
	return args.Error(0)
}

func logEventTestData() domain.LogEvent {
	return domain.LogEvent{
		Stream: domain.LogStdout,
		Line:   "The first line",
	}
}

func TestBuildLogService_AppendLog_Success(t *testing.T) {
	mockDB := new(mockBuildLogRepository)
	ctx := context.Background()
	buildLog := logEventTestData()
	buildId := "0"

	mockDB.On("Save", mock.Anything, mock.Anything).Return(nil)

	service := NewBuildLogService(mockDB)
	err := service.AppendLog(ctx, buildId, buildLog)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestBuildLogService_AppendLog_Error(t *testing.T) {
	mockDB := new(mockBuildLogRepository)
	expectedErr := errors.New("database error")
	ctx := context.Background()
	buildLog := logEventTestData()
	buildId := "0"

	mockDB.On("Save", mock.Anything, mock.Anything).Return(expectedErr)

	service := NewBuildLogService(mockDB)
	err := service.AppendLog(ctx, buildId, buildLog)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
