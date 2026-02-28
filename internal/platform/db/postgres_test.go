package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/H3nSte1n/ci-orchestrator/internal/testhelpers"
)

func TestNewPostgresConnection(t *testing.T) {
	cfg := testhelpers.GetTestConfig(t)

	db, err := NewPostgresConnection(cfg)
	require.NoError(t, err)
	require.NotNil(t, db)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	require.NoError(t, sqlDB.Ping())
}

func TestNewPostgresConnection_InvalidConfigs(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
	}{
		{
			name: "invalid host",
			setup: func() {
				cfg := testhelpers.GetTestConfig(t)
				cfg.DB.Host = "invalid-host-does-not-exist"
				_, err := NewPostgresConnection(cfg)
				assert.Error(t, err)
			},
		},
		{
			name: "invalid credentials",
			setup: func() {
				cfg := testhelpers.GetTestConfig(t)
				cfg.DB.Password = "wrong-password-12345"
				_, err := NewPostgresConnection(cfg)
				assert.Error(t, err)
			},
		},
		{
			name: "invalid database",
			setup: func() {
				cfg := testhelpers.GetTestConfig(t)
				cfg.DB.Name = "non_existent_database_xyz"
				_, err := NewPostgresConnection(cfg)
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
		})
	}
}
