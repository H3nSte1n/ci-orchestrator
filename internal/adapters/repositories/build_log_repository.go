package repositories

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"gorm.io/gorm"
)

type buildLogRepository struct {
	db ports.DB
}

func NewBuildLogRepository(db *gorm.DB) ports.BuildLogRepository {
	return &buildLogRepository{
		db: NewGormAdapter(db),
	}
}

func (blR *buildLogRepository) Save(ctx context.Context, buildLog *domain.BuildLog) error {
	return blR.db.WithContext(ctx).Create(buildLog).GetError()
}
