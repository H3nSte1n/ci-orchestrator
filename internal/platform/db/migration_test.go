package db

import (
	"testing"

	"github.com/H3nSte1n/ci-orchestrator/internal/testhelpers"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func TestMigrate(t *testing.T) {
	cfg := testhelpers.GetTestConfig(t)

	if err := Migrate(cfg); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}
}

func TestMigrateIdempotent(t *testing.T) {
	cfg := testhelpers.GetTestConfig(t)

	if err := Migrate(cfg); err != nil {
		t.Fatalf("First migration failed: %v", err)
	}

	if err := Migrate(cfg); err != nil {
		t.Fatalf("Second migration should not fail: %v", err)
	}
}

func TestMigrateInvalidConnection(t *testing.T) {
	cfg := testhelpers.GetTestConfig(t)

	cfg.DB.Host = "invalid-host-that-does-not-exist"

	err := Migrate(cfg)
	if err == nil {
		t.Error("Expected error for invalid host, got nil")
	}

	if err != nil && err.Error() == "" {
		t.Error("Error message should not be empty")
	}
}
