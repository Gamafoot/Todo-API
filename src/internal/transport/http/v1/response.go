package v1

import (
	"github.com/labstack/echo/v4"
)

type response struct {
	Message string `json:"message"`
}

func newResponse(c echo.Context, status int, message string) error {
	return c.JSON(status, response{message})
}
