package v1

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (h *handler) initSwaggerRoutes(api *echo.Group) {
	api.GET("/swagger*", echoSwagger.WrapHandler)
}
