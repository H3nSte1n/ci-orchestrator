package db

import (
	"testing"

	"github.com/H3nSte1n/ci-orchestrator/internal/testhelpers"
)

func TestNewPostgresConnection(t *testing.T) {
	cfg := testhelpers.GetTestConfig(t)

	db, err := NewPostgresConnection(cfg)
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}

	if db == nil {
		t.Error("Expected database connection, got nil")
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying SQL database: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestNewPostgresConnectionInvalidHost(t *testing.T) {
	cfg := testhelpers.GetTestConfig(t)

	cfg.DB.Host = "invalid-host-does-not-exist"

	db, err := NewPostgresConnection(cfg)
	if err == nil {
		t.Error("Expected error for invalid host, got nil")
	}

	if db != nil {
		t.Error("Expected nil database for failed connection")
	}
}

func TestNewPostgresConnectionInvalidCredentials(t *testing.T) {
	cfg := testhelpers.GetTestConfig(t)

	cfg.DB.Password = "wrong-password-12345"

	db, err := NewPostgresConnection(cfg)
	if err == nil {
		t.Error("Expected error for invalid credentials, got nil")
	}

	if db != nil {
		t.Error("Expected nil database for failed connection")
	}
}

func TestNewPostgresConnectionInvalidDatabase(t *testing.T) {
	cfg := testhelpers.GetTestConfig(t)

	cfg.DB.Name = "non_existent_database_xyz"

	db, err := NewPostgresConnection(cfg)
	if err == nil {
		t.Error("Expected error for non-existent database, got nil")
	}

	if db != nil {
		t.Error("Expected nil database for failed connection")
	}
}
