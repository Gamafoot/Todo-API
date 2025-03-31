package v1

import (
	"errors"
	"fmt"
	"net/http"
	"root/internal/domain"

	"github.com/labstack/echo/v4"
)

func (h *handler) initProjectRoutes(api *echo.Group) {
	api.POST("/projects", h.FindProjects)
	api.POST("/projects", h.CreateProject)
	api.PATCH("/projects/:project_id", h.UpdateProject)
	api.DELETE("/projects/:project_id", h.DeleteProject)
}

func (h *handler) FindProjects(c echo.Context) error {
	page, err := getIntFromQuery(c, "page")
	if err != nil {
		return newResponse(c, http.StatusBadRequest, "page is not digit")
	}

	limit, err := getIntFromQuery(c, "limit")
	if err != nil {
		return newResponse(c, http.StatusBadRequest, "limit is not digit")
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	columns, amount, err := h.service.Project.FindAll(userId, page, limit)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return newResponse(c, http.StatusForbidden, err.Error())
		}

		return err
	}

	c.Response().Header().Set("X-Total-Count", fmt.Sprintf("%d", amount))

	return c.JSON(http.StatusOK, columns)
}

func (h *handler) CreateProject(c echo.Context) error {
	input := new(domain.CreateProjectInput)

	if err := c.Bind(input); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	project, err := h.service.Project.Create(userId, input)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, project)
}

func (h *handler) UpdateProject(c echo.Context) error {
	input := new(domain.UpdateProjectInput)

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

	project, err := h.service.Project.Update(userId, projectId, input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return newResponse(c, http.StatusForbidden, err.Error())
		} else if errors.Is(err, domain.ErrRecordNotFound) {
			return newResponse(c, http.StatusNotFound, err.Error())
		}

		return err
	}

	return c.JSON(http.StatusOK, project)
}

func (h *handler) DeleteProject(c echo.Context) error {
	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	projectId, err := getUIntFromParam(c, "project_id")
	if err != nil {
		return err
	}

	err = h.service.Project.Delete(userId, projectId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return newResponse(c, http.StatusForbidden, err.Error())
		} else if errors.Is(err, domain.ErrRecordNotFound) {
			return newResponse(c, http.StatusNotFound, err.Error())
		}

		return err
	}

	return c.NoContent(http.StatusNoContent)
}
