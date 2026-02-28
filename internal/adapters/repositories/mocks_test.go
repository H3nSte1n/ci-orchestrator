package repositories

import (
	"context"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/ports"
	"github.com/stretchr/testify/mock"
)

type mockDB struct {
	mock.Mock
	Error error
}

func (m *mockDB) WithContext(ctx context.Context) ports.DB {
	m.Called(ctx)
	return m
}

func (m *mockDB) Create(value interface{}) ports.DB {
	m.Called(value)
	return m
}

func (m *mockDB) Where(query interface{}, args ...interface{}) ports.DB {
	m.Called(query, args)
	return m
}

func (m *mockDB) Updates(value interface{}) ports.DB {
	m.Called(value)
	return m
}

func (m *mockDB) Transaction(fn func(tx ports.DB) error) error {
	err := fn(m)
	if err != nil {
		return err
	}
	return m.Error
}

func (m *mockDB) First(value interface{}) ports.DB {
	m.Called(value)
	return m
}

func (m *mockDB) Order(value string) ports.DB {
	m.Called(value)
	return m
}

func (m *mockDB) Clauses(conds ...interface{}) ports.DB {
	m.Called(conds)
	return m
}

func (m *mockDB) Model(value interface{}) ports.DB {
	m.Called(value)
	return m
}

func (m *mockDB) GetError() error {
	return m.Error
}
