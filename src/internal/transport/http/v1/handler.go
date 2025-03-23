package v1

import (
	"fmt"
	"root/internal/config"
	"root/internal/service"
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
	fmt.Printf("aaa: %v\n", h.config.Cors.AllowOrigins)
	api.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     h.config.Cors.AllowOrigins,
			AllowMethods:     h.config.Cors.AllowMethods,
			AllowHeaders:     h.config.Cors.AllowHeaders,
			AllowCredentials: true,
		}),
		middleware.Logger(),
	)

	v1 := api.Group("/v1")
	{
		h.initAuthRoutes(v1)

		requiredAuth := v1.Group("", h.requiredAuth)
		{
			h.initTaskRoutes(requiredAuth)
			h.initProjectRoutes(requiredAuth)

			projectsId := requiredAuth.Group("/projects/:id")
			{
				h.initColumnRoutes(projectsId)
			}
		}
	}
}
