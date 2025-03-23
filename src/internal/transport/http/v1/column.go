package v1

import (
	"errors"
	"fmt"
	"net/http"
	"root/internal/domain"
	"strconv"

	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

func (h *handler) initColumnRoutes(api *echo.Group) {
	api.GET("/columns", h.CreateColumn)
	api.POST("/columns", h.CreateColumn)
	api.PATCH("/columns", h.UpdateProject)
	api.DELETE("/columns", h.DeleteProject)
}

func (h *handler) FindColumns(c echo.Context) error {
	page, err := getIntFromParam(c, "page")
	if err != nil {
		return newResponse(c, http.StatusBadRequest, "page is not digit")
	}

	limit, err := getIntFromParam(c, "limit")
	if err != nil {
		return newResponse(c, http.StatusBadRequest, "limit is not digit")
	}

	projectId, err := getProjectId(c)
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

type CreateColumnInput struct {
	Name string `json:"username" binding:"required,min=3,max=50"`
}

func (h *handler) CreateColumn(c echo.Context) error {
	input := new(CreateColumnInput)

	if err := c.Bind(input); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	projectId, err := getProjectId(c)
	if err != nil {
		return err
	}

	err = h.service.Column.Create(userId, projectId, input.Name)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return newResponse(c, http.StatusBadRequest, err.Error())
		}

		return err
	}

	return c.NoContent(http.StatusCreated)
}

type UpdateColumnInput struct {
	Id   uint   `json:"id" binding:"required"`
	Name string `json:"username" binding:"required,min=3,max=50"`
}

func (h *handler) UpdateColumn(c echo.Context) error {
	input := new(UpdateProjectInput)

	if err := c.Bind(input); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	err = h.service.Project.Update(userId, input.Id, input.Name)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) || errors.Is(err, domain.ErrRecordNotFound) {
			return newResponse(c, http.StatusBadRequest, err.Error())
		}

		return err
	}

	return c.NoContent(http.StatusOK)
}

type DeleteColumnInput struct {
	Id uint `json:"id" binding:"required"`
}

func (h *handler) DeleteColumn(c echo.Context) error {
	input := new(DeleteProjectInput)

	if err := c.Bind(input); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	err = h.service.Project.Delete(userId, input.Id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) || errors.Is(err, domain.ErrRecordNotFound) {
			return newResponse(c, http.StatusBadRequest, err.Error())
		}

		return err
	}

	return c.NoContent(http.StatusOK)
}

func getProjectId(c echo.Context) (uint, error) {
	projectId := c.Param("id")

	projectIdInt, err := strconv.Atoi(projectId)
	if err != nil {
		return 0, pkgErrors.WithStack(err)
	}

	return uint(projectIdInt), nil
}
