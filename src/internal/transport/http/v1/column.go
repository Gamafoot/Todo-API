package v1

import (
	"errors"
	"fmt"
	"net/http"
	"root/internal/domain"

	"github.com/labstack/echo/v4"
)

func (h *handler) initColumnRoutes(api *echo.Group) {
	api.GET("/projects/:project_id/columns", h.FindColumns)
	api.POST("/columns", h.CreateColumn)
	api.PATCH("/columns/:column_id", h.UpdateColumn)
	api.DELETE("/columns/:column_id", h.DeleteColumn)
}

// @Summary Список колонок
// @Tags column
// @Produce json
// @Security BearerAuth
// @Param project_id path int true "ID проекта"
// @Param page query int false "Номер страницы, по уполчанию 1"
// @Param limit query int false "Кол-во итоговых записей, по уполчанию 10"
// @Success 200 {array} domain.Column
// @Header 200 {integer} X-Total-Pages "Общее количество страниц колонок на проекте"
// @Failure 400
// @Failure 401
// @Router /api/v1/projects/{project_id}/columns [get]
func (h *handler) FindColumns(c echo.Context) error {
	page, err := getIntFromQuery(c, "page", 1)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	limit, err := getIntFromQuery(c, "limit", 10)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, err.Error())
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
			return c.NoContent(http.StatusForbidden)
		}

		return err
	}

	c.Response().Header().Set("Access-Control-Expose-Headers", "X-Total-Pages")
	c.Response().Header().Set("X-Total-Pages", fmt.Sprintf("%d", amount))

	return c.JSON(http.StatusOK, columns)
}

// @Summary Создать колонку
// @Tags column
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body domain.CreateColumnInput true "Данные для создания колонки"
// @Success 200 {object} domain.Column
// @Failure 400
// @Failure 401
// @Router /api/v1/columns [post]
func (h *handler) CreateColumn(c echo.Context) error {
	input := new(domain.CreateColumnInput)

	if err := c.Bind(input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	column, err := h.service.Column.Create(userId, input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return c.NoContent(http.StatusForbidden)
		}

		return err
	}

	return c.JSON(http.StatusCreated, column)
}

// @Summary Обновить колонку
// @Tags column
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param column_id path int true "ID колонки"
// @Param body body domain.UpdateColumnInput true "Данные для обновления колонки"
// @Success 200 {object} domain.Column
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /api/v1/columns/{column_id} [patch]
func (h *handler) UpdateColumn(c echo.Context) error {
	input := new(domain.UpdateColumnInput)

	if err := c.Bind(input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	columnId, err := getUIntFromParam(c, "column_id")
	if err != nil {
		return err
	}

	column, err := h.service.Column.Update(userId, columnId, input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return c.NoContent(http.StatusForbidden)
		} else if errors.Is(err, domain.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return err
	}

	return c.JSON(http.StatusOK, column)
}

// @Summary Удалить колонку
// @Tags column
// @Produce json
// @Security BearerAuth
// @Param column_id path int true "ID колонки"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /api/v1/columns/{column_id} [delete]
func (h *handler) DeleteColumn(c echo.Context) error {
	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	columnId, err := getUIntFromParam(c, "column_id")
	if err != nil {
		return err
	}

	err = h.service.Column.Delete(userId, columnId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return c.NoContent(http.StatusForbidden)
		}

		return err
	}

	return c.NoContent(http.StatusNoContent)
}
