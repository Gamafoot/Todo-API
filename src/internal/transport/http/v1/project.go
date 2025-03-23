package v1

import (
	"errors"
	"fmt"
	"net/http"
	"root/internal/domain"

	"github.com/labstack/echo/v4"
)

func (h *handler) initProjectRoutes(api *echo.Group) {
	api.POST("/projects", h.CreateProject)
	api.PATCH("/projects", h.UpdateProject)
	api.DELETE("/projects", h.DeleteProject)
}

func (h *handler) FindProjects(c echo.Context) error {
	page, err := getIntFromParam(c, "page")
	if err != nil {
		return newResponse(c, http.StatusBadRequest, "page is not digit")
	}

	limit, err := getIntFromParam(c, "limit")
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
			return newResponse(c, 400, err.Error())
		}

		return err
	}

	c.Response().Header().Set("X-Total-Count", fmt.Sprintf("%d", amount))

	return c.JSON(http.StatusOK, columns)
}

type CreateProjectInput struct {
	Name string `json:"username" binding:"required,min=3,max=50"`
}

func (h *handler) CreateProject(c echo.Context) error {
	input := new(CreateProjectInput)

	if err := c.Bind(input); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	err = h.service.Project.Create(userId, input.Name)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

type UpdateProjectInput struct {
	Id   uint   `json:"id" binding:"required"`
	Name string `json:"username" binding:"required,min=3,max=50"`
}

func (h *handler) UpdateProject(c echo.Context) error {
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
			return newResponse(c, 400, err.Error())
		}

		return err
	}

	return c.NoContent(http.StatusOK)
}

type DeleteProjectInput struct {
	Id uint `json:"id" binding:"required"`
}

func (h *handler) DeleteProject(c echo.Context) error {
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
			return newResponse(c, 400, err.Error())
		}

		return err
	}

	return c.NoContent(http.StatusOK)
}
