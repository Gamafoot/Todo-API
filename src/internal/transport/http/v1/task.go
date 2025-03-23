package v1

import (
	"fmt"
	"net/http"
	"root/internal/domain"
	"root/internal/service"
	"strconv"

	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

func (h *handler) initTaskRoutes(api *echo.Group) {
	api.GET("/tasks", h.GetTasks)
	api.POST("/tasks", h.SaveTask)
	api.DELETE("/tasks/:taskId", h.Delete)
}

func (h *handler) GetTasks(c echo.Context) error {
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return newResponse(c, http.StatusBadRequest, "page is not digit")
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return newResponse(c, http.StatusBadRequest, "limit is not digit")
	}

	if limitInt > 10 {
		limitInt = 10
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	tasks, amount, err := h.service.Task.GetChunk(userId, pageInt, limitInt)
	if err != nil {
		if pkgErrors.Is(err, domain.ErrTaskNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return err
	}

	c.Response().Header().Set("X-Total-Count", fmt.Sprintf("%d", amount))

	return c.JSON(http.StatusOK, tasks)
}

type requestTaskInput struct {
	Name        string `json:"name" binding:"required,min=1,max=60"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Deadline    string `json:"deadline"`
}

func (h *handler) SaveTask(c echo.Context) error {
	inp := requestTaskInput{}

	if err := c.Bind(&inp); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	err = h.service.Task.Save(&service.TaskInput{
		UserId:      userId,
		Name:        inp.Name,
		Description: inp.Description,
		Status:      inp.Status,
		Deadline:    inp.Deadline,
	})
	if err != nil {
		if pkgErrors.Is(err, domain.ErrInvalidDeadlineFormat) {
			return newResponse(c, http.StatusBadRequest, err.Error())
		}
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) Delete(c echo.Context) error {
	taskId := c.Param("taskId")

	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		return newResponse(c, http.StatusBadRequest, "taskId is not digit")
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	if err = h.service.Task.Delete(uint(taskIdInt), userId); err != nil {
		if pkgErrors.Is(err, domain.ErrTaskNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}

	return c.NoContent(http.StatusOK)
}
