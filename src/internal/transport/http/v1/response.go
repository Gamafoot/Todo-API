package v1

import (
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(c echo.Context, status int, message string) error {
	return c.JSON(status, ErrorResponse{message})
}
