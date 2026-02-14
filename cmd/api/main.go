package main

import (
	"github.com/H3nSte1n/ci-orchestrator/internal/platform/config"
	"github.com/H3nSte1n/ci-orchestrator/internal/platform/db"
	"os"
)

func main() {
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	cfg, err := config.LoadConfig(env)
	if err != nil {
		panic(err)
	}

	if db.Migrate(cfg) != nil {
		panic(err)
	}
}
