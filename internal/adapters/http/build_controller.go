package http

import (
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BuildController struct {
	buildService ports.BuildService
}

func NewBuildController(buildService ports.BuildService) *BuildController {
	return &BuildController{
		buildService: buildService,
	}
}

func (bc *BuildController) CreateBuild(c *gin.Context) {
	var build domain.Build

	if err := c.ShouldBindJSON(&build); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if build.RepoUrl == "" || build.Ref == "" || build.Command == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields: repo_url, ref, command"})
		return
	}

	if err := bc.buildService.CreateBuild(c.Request.Context(), &build); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create build", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, build)
}

func (bc *BuildController) CancelBuild(c *gin.Context) {
	buildId := c.Param("id")

	if buildId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "build id is required"})
		return
	}

	if err := bc.buildService.CancelBuild(c.Request.Context(), buildId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel build", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Build canceled successfully"})
}

func (bc *BuildController) UpdateStatus(c *gin.Context) {
	buildId := c.Param("id")

	if buildId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "build id is required"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	status := domain.BuildStatus(req.Status)

	switch status {
	case domain.BuildStatusPending, domain.BuildStatusRunning, domain.BuildStatusSuccess, domain.BuildStatusFailed, domain.BuildStatusCanceled:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status", "valid_statuses": []string{
			string(domain.BuildStatusPending),
			string(domain.BuildStatusRunning),
			string(domain.BuildStatusSuccess),
			string(domain.BuildStatusFailed),
			string(domain.BuildStatusCanceled),
		}})
		return
	}

	if err := bc.buildService.UpdateStatus(c.Request.Context(), buildId, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Build status updated successfully"})
}

func (bc *BuildController) GetBuild(c *gin.Context) {
	buildId := c.Param("id")

	if buildId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "build id is required"})
		return
	}

	build, err := bc.buildService.GetBuild(c.Request.Context(), buildId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get build", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, build)
}
