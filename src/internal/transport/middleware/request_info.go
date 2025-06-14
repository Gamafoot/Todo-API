package middleware

import (
	"log"

	"github.com/labstack/echo/v4"
)

func RequestInfo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Printf("Method: %s\nPath: %s\nHeaders: %#v\n", c.Request().Method, c.Request().RequestURI, c.Request().Header)
		return next(c)
	}
}
