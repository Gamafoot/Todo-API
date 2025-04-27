package v1

import (
	"errors"
	"fmt"
	"net/http"
	"root/internal/domain"

	"github.com/labstack/echo/v4"
)

func (h *handler) initProjectRoutes(api *echo.Group) {
	api.GET("/projects", h.ListProjects)
	api.GET("/projects/:project_id", h.DetailProject)
	api.POST("/projects", h.CreateProject)
	api.PATCH("/projects/:project_id", h.UpdateProject)
	api.DELETE("/projects/:project_id", h.DeleteProject)
	api.GET("/projects/:project_id/stats", h.ProjectStats)
}

// @Summary Список проектов
// @Tags project
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы, по уполчанию 1"
// @Param limit query int false "Кол-во итоговых записей, по уполчанию 10"
// @Success 200 {array} domain.Project
// @Header 200 {integer} X-Total-Pages "Общее количество страниц проектов у пользователя"
// @Failure 400
// @Failure 401
// @Failure 403
// @Router /api/v1/projects [get]
func (h *handler) ListProjects(c echo.Context) error {
	page, err := getIntFromQuery(c, "page", 1)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	limit, err := getIntFromQuery(c, "limit", 10)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	columns, amount, err := h.service.Project.List(userId, page, limit)
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

// @Summary Детали проекта
// @Tags project
// @Produce json
// @Security BearerAuth
// @Param project_id path int true "ID проекта"
// @Success 200 {object} domain.Project
// @Failure 400
// @Failure 401
// @Failure 403
// @Router /api/v1/projects/{project_id} [get]
func (h *handler) DetailProject(c echo.Context) error {
	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	projectId, err := getUIntFromParam(c, "project_id")
	if err != nil {
		return err
	}

	project, err := h.service.Project.Detail(userId, projectId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return c.NoContent(http.StatusForbidden)
		}
		return err
	}

	return c.JSON(http.StatusOK, project)
}

// @Summary Создать проект
// @Tags project
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body domain.CreateProjectInput true "Данные для создания проекта"
// @Success 200 {object} domain.Project
// @Failure 400
// @Failure 401
// @Router /api/v1/projects [post]
func (h *handler) CreateProject(c echo.Context) error {
	input := new(domain.CreateProjectInput)

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

	project, err := h.service.Project.Create(userId, input)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, project)
}

// @Summary Обновить проект
// @Tags project
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path int true "ID проекта"
// @Param body body domain.UpdateProjectInput true "Данные для обновления проекта"
// @Success 200 {object} domain.Project
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /api/v1/projects/{project_id} [patch]
func (h *handler) UpdateProject(c echo.Context) error {
	input := new(domain.UpdateProjectInput)

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

	projectId, err := getUIntFromParam(c, "project_id")
	if err != nil {
		return err
	}

	project, err := h.service.Project.Update(userId, projectId, input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return c.NoContent(http.StatusForbidden)
		} else if errors.Is(err, domain.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}

	return c.JSON(http.StatusOK, project)
}

// @Summary Удалить проект
// @Tags project
// @Produce json
// @Security BearerAuth
// @Param project_id path int true "ID проекта"
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /api/v1/projects/{project_id} [delete]
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
			return c.NoContent(http.StatusForbidden)
		}
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary Статистика проекта
// @Tags project
// @Produce json
// @Security BearerAuth
// @Param project_id path int true "ID проекта"
// @Success 200 {object} domain.ProjectStats
// @Failure 401
// @Failure 404
// @Router /api/v1/projects/{project_id}/stats [get]
func (h *handler) ProjectStats(c echo.Context) error {
	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	projectId, err := getUIntFromParam(c, "project_id")
	if err != nil {
		return err
	}

	stats, err := h.service.Project.GetStats(userId, projectId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}

	return c.JSON(http.StatusOK, stats)
}
