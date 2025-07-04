package v1

import (
	"root/internal/config"
	"root/internal/service"
	mymiddleware "root/internal/transport/middleware"
	"root/pkg/jwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type handler struct {
	config       *config.Config
	service      *service.Service
	tokenManager jwt.TokenManager
}

func NewHandler(cfg *config.Config, service *service.Service, tokenManager jwt.TokenManager) *handler {
	return &handler{
		config:       cfg,
		service:      service,
		tokenManager: tokenManager,
	}
}

func (h *handler) InitRoutes(api *echo.Group) {
	api.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     h.config.Cors.AllowOrigins,
			AllowMethods:     h.config.Cors.AllowMethods,
			AllowHeaders:     h.config.Cors.AllowHeaders,
			AllowCredentials: true,
		}),
		mymiddleware.Logger,
	)

	h.initSwaggerRoutes(api)

	v1 := api.Group("/api/v1")
	{
		h.initAuthRoutes(v1)

		requiredAuth := v1.Group("", h.requiredAuth)
		{
			h.initProjectRoutes(requiredAuth)
			h.initColumnRoutes(requiredAuth)
			h.initTaskRoutes(requiredAuth)
			h.initSubTaskRoutes(requiredAuth)
			h.initStatsRoutes(requiredAuth)
			h.initHeatmapRoutes(requiredAuth)
		}
	}
}
