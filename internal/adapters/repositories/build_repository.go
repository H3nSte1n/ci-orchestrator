package repositories

import (
	"context"
	"errors"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/domain"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
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

func (r *buildRepository) ClaimNext(ctx context.Context, workerId string) (*domain.Build, error) {
	var build domain.Build

	err := r.db.WithContext(ctx).Transaction(func(tx ports.DB) error {
		if err := tx.Where("status = ?", domain.BuildStatusPending).
			Where("locked_by IS NULL").
			Where("locked_at IS NULL").
			Order("created_at ASC").
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			First(&build).GetError(); err != nil {
			return err
		}

		if err := tx.Model(&build).Updates(map[string]interface{}{
			"status":    domain.BuildStatusRunning,
			"locked_by": workerId,
			"locked_at": time.Now(),
		}).GetError(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &build, nil
}
