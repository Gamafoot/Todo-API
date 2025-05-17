package v1

import (
	"math"
	"root/internal/domain"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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

func getIntFromQuery(c echo.Context, param string, defaultValue ...int) (int, error) {
	value := c.QueryParam(param)

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return 0, errors.Wrap(domain.ErrNotDigit, param)
	}

	if valueInt < 0 {
		return 0, errors.Wrap(domain.ErrNotPositiveDigit, param)
	}

	return valueInt, nil
}

func getUIntFromParam(c echo.Context, param string) (uint, error) {
	value := c.Param(param)

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.Wrap(domain.ErrNotDigit, param)
	}

	return uint(valueInt), nil
}

func getBoolFromQuery(c echo.Context, param string) (bool, error) {
	value := c.QueryParam(param)

	if value == "true" {
		return true, nil
	} else if value == "false" || len(value) == 0 {
		return false, nil
	}

	return false, errors.New("wrong value for boolean")
}

func getStringFromQuery(c echo.Context, param string, defaultValue string) (bool, error) {
	value := c.QueryParam(param)

	if value == "true" {
		return true, nil
	} else if value == "false" || len(value) == 0 {
		return false, nil
	}

	return false, errors.New("wrong value for boolean")
}

func getPageCount(count, limit int) int {
	return int(math.Ceil(float64(count) / float64(limit)))
}
