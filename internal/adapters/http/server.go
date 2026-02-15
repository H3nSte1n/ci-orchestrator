package http

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine     *gin.Engine
	controller *BuildController
}

func NewRouter(controller *BuildController) *Router {
	engine := gin.Default()

	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	return &Router{
		engine:     engine,
		controller: controller,
	}
}

func (r *Router) RegisterRoutes() {
	v1 := r.engine.Group("/api/v1")
	{
		builds := v1.Group("/builds")
		{
			builds.POST("", r.controller.CreateBuild)
			builds.GET("/:id", r.controller.GetBuild)
			builds.PATCH("/:id/status", r.controller.UpdateStatus)
			builds.POST("/:id/cancel", r.controller.CancelBuild)
		}
	}
}

func (r *Router) Run(addr string) error {
	r.RegisterRoutes()
	return r.engine.Run(addr)
}
