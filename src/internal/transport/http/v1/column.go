package v1

import (
	"errors"
	"fmt"
	"net/http"
	"root/internal/domain"

	"github.com/labstack/echo/v4"
)

func (h *handler) initColumnRoutes(api *echo.Group) {
	api.GET("/columns", h.FindColumns)
	api.POST("/columns", h.CreateColumn)
	api.PATCH("/columns/:column_id", h.UpdateProject)
	api.DELETE("/columns/:column_id", h.DeleteProject)
}

func (h *handler) FindColumns(c echo.Context) error {
	page, err := getIntFromQuery(c, "page")
	if err != nil {
		return newResponse(c, http.StatusBadRequest, "page is not digit")
	}

	limit, err := getIntFromQuery(c, "limit")
	if err != nil {
		return newResponse(c, http.StatusBadRequest, "limit is not digit")
	}

	projectId, err := getUIntFromParam(c, "project_id")
	if err != nil {
		return err
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	columns, amount, err := h.service.Column.FindAll(userId, projectId, page, limit)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return newResponse(c, 400, err.Error())
		}

		return err
	}

	c.Response().Header().Set("X-Total-Count", fmt.Sprintf("%d", amount))

	return c.JSON(http.StatusOK, columns)
}

func (h *handler) CreateColumn(c echo.Context) error {
	input := new(domain.CreateColumnInput)

	if err := c.Bind(input); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	column, err := h.service.Column.Create(userId, input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return newResponse(c, http.StatusBadRequest, err.Error())
		}

		return err
	}

	return c.JSON(http.StatusCreated, column)
}

func (h *handler) UpdateColumn(c echo.Context) error {
	input := new(domain.UpdateColumnInput)

	if err := c.Bind(input); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	projectId, err := getUIntFromParam(c, "project_id")
	if err != nil {
		return err
	}

	column, err := h.service.Column.Update(userId, projectId, input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return newResponse(c, http.StatusForbidden, err.Error())
		} else if errors.Is(err, domain.ErrRecordNotFound) {
			return newResponse(c, http.StatusNotFound, err.Error())
		}

		return err
	}

	return c.JSON(http.StatusOK, column)
}

func (h *handler) DeleteColumn(c echo.Context) error {
	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	columnId, err := getUIntFromParam(c, "column_id")
	if err != nil {
		return err
	}

	err = h.service.Project.Delete(userId, columnId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return newResponse(c, http.StatusForbidden, err.Error())
		} else if errors.Is(err, domain.ErrRecordNotFound) {
			return newResponse(c, http.StatusNotFound, err.Error())
		}

		return err
	}

	return c.NoContent(http.StatusOK)
}
