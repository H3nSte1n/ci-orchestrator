package ports

import (
	"context"
)

type DB interface {
	WithContext(ctx context.Context) DB
	Create(value interface{}) DB
	Where(query interface{}, args ...interface{}) DB
	Updates(value interface{}) DB
	First(value interface{}) DB
	GetError() error
	Transaction(f func(tx DB) error) error
	Order(value string) DB
	Clauses(conds ...interface{}) DB
	Model(value interface{}) DB
}
