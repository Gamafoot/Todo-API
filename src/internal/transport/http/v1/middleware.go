package v1

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	userCtx = "userId"
)

func (h *handler) requiredAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")

		userId, err := h.tokenManager.Parse(token)
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		c.Set(userCtx, userId)

		return next(c)
	}
}
