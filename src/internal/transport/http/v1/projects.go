package v1

import (
	"errors"
	"net/http"
	"root/internal/domain"
	"root/internal/service"
	customErrors "root/pkg/errors"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *handler) initProjectRoutes(api *echo.Group) {
	api.POST("/auth/login", h.Login)
	api.POST("/auth/register", h.Register)
	api.GET("/auth/refresh", h.RefreshToken, h.requiredAuth)
}

func (h *handler) CreateProject(c echo.Context) error {
	input := new(loginInput)

	if err := c.Bind(input); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	tokens, err := h.service.Auth.Login(service.LoginInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidLoginOrPassword) {
			return newResponse(c, http.StatusUnauthorized, err.Error())
		}
		return err
	}

	return setTokensToResponse(c, tokens, h.config.Jwt.RefreshTokenTtl)
}
