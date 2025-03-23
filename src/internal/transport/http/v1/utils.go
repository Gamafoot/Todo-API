package v1

import (
	"strconv"

	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

func getUserIdFromContext(c echo.Context) (uint, error) {
	value := c.Get(userCtx)

	userId, ok := value.(uint)
	if !ok {
		return 0, pkgErrors.New("userCtx is of invalid type")
	}

	return userId, nil
}

func getIntFromParam(c echo.Context, param string) (int, error) {
	value := c.QueryParam(param)

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return valueInt, nil
}
