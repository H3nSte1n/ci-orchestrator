package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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

func TestBuildController_CreateBuild_Success(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("CreateBuild", mock.Anything, mock.MatchedBy(func(b *domain.Build) bool {
		return b.RepoUrl == "https://github.com/test/repo" &&
			b.Ref == "main" &&
			b.Command == "npm test"
	})).Return(nil)
	bc := NewBuildController(mockBuildService)

	router := gin.New()
	router.POST("/builds", bc.CreateBuild)

	body := []byte(`{"repo_url": "https://github.com/test/repo","ref": "main","command": "npm test"}`)
	req := httptest.NewRequest("POST", "/builds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var got domain.Build
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/test/repo", got.RepoUrl)
	assert.Equal(t, "main", got.Ref)
	assert.Equal(t, "npm test", got.Command)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockBuildService.AssertExpectations(t)
}

func TestBuildController_CreateBuild_Error(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("CreateBuild", mock.Anything, mock.Anything).Return(assert.AnError)

	bc := NewBuildController(mockBuildService)

	router := gin.New()
	router.POST("/builds", bc.CreateBuild)

	body := []byte(`{"repo_url": "https://github.com/test/repo","ref": "main","command": "npm test"}`)
	req := httptest.NewRequest("POST", "/builds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Contains(t, w.Body.String(), "failed to create build")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockBuildService.AssertExpectations(t)
}

func TestBuildController_CreateBuild_MissingRequiredFields(t *testing.T) {
	mockBuildService := new(mockBuildService)

	bc := NewBuildController(mockBuildService)

	router := gin.New()
	router.POST("/builds", bc.CreateBuild)

	body := []byte(`{"repo_url": "https://github.com/test/repo"}`)
	req := httptest.NewRequest("POST", "/builds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing required fields: repo_url, ref, command")
}

func TestBuildController_CreateBuild_InvalidJSON(t *testing.T) {
	mockBuildService := new(mockBuildService)

	bc := NewBuildController(mockBuildService)

	router := gin.New()
	router.POST("/builds", bc.CreateBuild)

	body := []byte(`{"repo_url": invalid json}`)
	req := httptest.NewRequest("POST", "/builds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}

func TestBuildController_CancelBuild_Success(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("CancelBuild", mock.Anything, "test-id").Return(nil)

	bc := NewBuildController(mockBuildService)

	router := gin.New()
	router.POST("/builds/:id/cancel", bc.CancelBuild)

	req := httptest.NewRequest("POST", "/builds/test-id/cancel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Contains(t, w.Body.String(), "Build canceled successfully")
	assert.Equal(t, http.StatusOK, w.Code)
	mockBuildService.AssertExpectations(t)
}

func TestBuildController_CancelBuild_Error(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("CancelBuild", mock.Anything, "test-id").Return(assert.AnError)

	bc := NewBuildController(mockBuildService)

	router := gin.New()
	router.POST("/builds/:id/cancel", bc.CancelBuild)

	req := httptest.NewRequest("POST", "/builds/test-id/cancel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockBuildService.AssertExpectations(t)
}

func TestBuildController_GetBuild_Success(t *testing.T) {
	expectedBuild := &domain.Build{
		ID:      "test-id",
		RepoUrl: "https://github.com/test/repo",
		Ref:     "main",
		Command: "npm test",
	}
	mockBuildService := new(mockBuildService)
	mockBuildService.On("GetBuild", mock.Anything, "test-id").Return(expectedBuild, nil)

	bc := NewBuildController(mockBuildService)

	router := gin.New()
	router.GET("/builds/:id", bc.GetBuild)

	req := httptest.NewRequest("GET", "/builds/test-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var returnedBuild domain.Build
	err := json.Unmarshal(w.Body.Bytes(), &returnedBuild)
	assert.NoError(t, err)
	assert.Equal(t, expectedBuild.ID, returnedBuild.ID)
	assert.Equal(t, expectedBuild.RepoUrl, returnedBuild.RepoUrl)
	assert.Equal(t, expectedBuild.Ref, returnedBuild.Ref)
	assert.Equal(t, expectedBuild.Command, returnedBuild.Command)

	mockBuildService.AssertExpectations(t)
}

func TestBuildController_GetBuild_Error(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("GetBuild", mock.Anything, "test-id").Return(nil, assert.AnError)

	bc := NewBuildController(mockBuildService)

	router := gin.New()
	router.GET("/builds/:id", bc.GetBuild)

	req := httptest.NewRequest("GET", "/builds/test-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockBuildService.AssertExpectations(t)
}

func TestBuildController_UpdateStatus_Success(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("UpdateStatus", mock.Anything, "test-id", domain.BuildStatusRunning).Return(nil)

	bc := NewBuildController(mockBuildService)

	router := gin.New()
	router.PATCH("/builds/:id/status", bc.UpdateStatus)

	body := []byte(`{"status": "running"}`)
	req := httptest.NewRequest("PATCH", "/builds/test-id/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBuildService.AssertExpectations(t)
}

func TestBuildController_UpdateStatus_Error(t *testing.T) {
	mockBuildService := new(mockBuildService)
	mockBuildService.On("UpdateStatus", mock.Anything, "test-id", domain.BuildStatusRunning).Return(assert.AnError)

	bc := NewBuildController(mockBuildService)

	router := gin.New()
	router.PATCH("/builds/:id/status", bc.UpdateStatus)

	body := []byte(`{"status": "running"}`)
	req := httptest.NewRequest("PATCH", "/builds/test-id/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockBuildService.AssertExpectations(t)
}

func TestBuildController_UpdateStatus_InvalidStatus(t *testing.T) {
	mockBuildService := new(mockBuildService)

	bc := NewBuildController(mockBuildService)

	router := gin.New()
	router.PATCH("/builds/:id/status", bc.UpdateStatus)

	body := []byte(`{"status": "invalid_status"}`)
	req := httptest.NewRequest("PATCH", "/builds/test-id/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid status")
	assert.Contains(t, w.Body.String(), "valid_statuses")
}

func TestBuildController_UpdateStatus_MissingStatusField(t *testing.T) {
	mockBuildService := new(mockBuildService)

	bc := NewBuildController(mockBuildService)

	router := gin.New()
	router.PATCH("/builds/:id/status", bc.UpdateStatus)

	body := []byte(`{}`)
	req := httptest.NewRequest("PATCH", "/builds/test-id/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}
