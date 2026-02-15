package main

import (
	"github.com/H3nSte1n/ci-orchestrator/internal/adapters/http"
	"github.com/H3nSte1n/ci-orchestrator/internal/adapters/repositories"
	"github.com/H3nSte1n/ci-orchestrator/internal/core/service"
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

	if err := db.Migrate(cfg); err != nil {
		panic(err)
	}

	dbConnection, err := db.NewPostgresConnection(cfg)
	if err != nil {
		panic(err)
	}

	buildRepository := repositories.NewBuildRepository(dbConnection)
	buildService := service.NewBuildService(buildRepository)
	buildController := http.NewBuildController(buildService)
	router := http.NewRouter(buildController)

	if err := router.Run(":" + cfg.ApiServiceConfig.Port); err != nil {
		panic(err)
	}
}
