package testhelpers

import (
	"os"
	"testing"

	"github.com/H3nSte1n/ci-orchestrator/internal/platform/config"
)

func GetTestConfig(t *testing.T) *config.Config {
	if os.Getenv("BASE_DIR") == "" {
		os.Setenv("BASE_DIR", "../../../")
	}

	cfg, err := config.LoadConfig("test")
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	return cfg
}
