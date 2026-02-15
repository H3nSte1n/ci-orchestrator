package db

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"gorm.io/gorm"
)

type buildRepository struct {
	db ports.DB
}

func NewBuildRepository(gormDB *gorm.DB) ports.BuildRepository {
	return &buildRepository{
		db: NewGormAdapter(gormDB),
	}
}

func (r *buildRepository) Save(ctx context.Context, build *domain.Build) error {
	return r.db.WithContext(ctx).Create(build).GetError()
}

func (r *buildRepository) Update(ctx context.Context, build *domain.Build) error {
	return r.db.WithContext(ctx).Where("id = ?", build.ID).Updates(build).GetError()
}

func (r *buildRepository) FindByID(ctx context.Context, buildId string) (*domain.Build, error) {
	var build domain.Build
	err := r.db.WithContext(ctx).Where("id = ?", buildId).First(&build).GetError()
	if err != nil {
		return nil, err
	}
	return &build, nil
}
