package db

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"gorm.io/gorm"
)

type gormAdapter struct {
	*gorm.DB
}

func NewGormAdapter(db *gorm.DB) ports.DB {
	return &gormAdapter{db}
}

func (g *gormAdapter) WithContext(ctx context.Context) ports.DB {
	return &gormAdapter{g.DB.WithContext(ctx)}
}

func (g *gormAdapter) Create(value interface{}) ports.DB {
	return &gormAdapter{g.DB.Create(value)}
}

func (g *gormAdapter) Where(query interface{}, args ...interface{}) ports.DB {
	return &gormAdapter{g.DB.Where(query, args...)}
}

func (g *gormAdapter) Updates(value interface{}) ports.DB {
	return &gormAdapter{g.DB.Updates(value)}
}

func (g *gormAdapter) First(value interface{}) ports.DB {
	return &gormAdapter{g.DB.First(value)}
}

func (g *gormAdapter) GetError() error {
	return g.DB.Error
}
