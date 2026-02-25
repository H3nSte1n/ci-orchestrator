package repositories

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (g *gormAdapter) Order(value string) ports.DB {
	return &gormAdapter{g.DB.Order(value)}
}

func (g *gormAdapter) Clauses(conds ...interface{}) ports.DB {
	clauseExprs := make([]clause.Expression, len(conds))
	for i, cond := range conds {
		clauseExprs[i] = cond.(clause.Expression)
	}
	return &gormAdapter{g.DB.Clauses(clauseExprs...)}
}

func (g *gormAdapter) Model(value interface{}) ports.DB {
	return &gormAdapter{g.DB.Model(value)}
}

func (g *gormAdapter) Transaction(fn func(tx ports.DB) error) error {
	return g.DB.Transaction(func(tx *gorm.DB) error {
		return fn(&gormAdapter{tx})
	})
}

func (g *gormAdapter) GetError() error {
	return g.DB.Error
}
