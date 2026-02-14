package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("BASE_DIR", "../../../")

	cfg, err := LoadConfig("test")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.App != (AppConfig{
		Name: "ci orchestrator",
		Env:  "test",
	}) {
		t.Errorf("App config mismatch. Got: %+v", cfg.App)
	}

	if cfg.DB != (DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "your_password",
		Name:     "ci_db_test",
		SSLMode:  "disable",
	}) {
		t.Errorf("DB config mismatch. Got: %+v", cfg.DB)
	}
}

func TestLoadConfigNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent config, got nil")
	}
}
