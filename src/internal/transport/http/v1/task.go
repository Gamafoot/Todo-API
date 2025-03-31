package v1

import (
	"errors"
	"fmt"
	"net/http"
	"root/internal/domain"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) initTaskRoutes(api *echo.Group) {
	api.GET("/tasks", h.GetTasks)
	api.POST("/tasks", h.CreateTask)
	api.PATCH("/tasks/:task_id", h.UpdateTask)
	api.DELETE("/tasks/:task_id", h.Delete)
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

	columnId, err := getUIntFromParam(c, "column_id")
	if err != nil {
		return err
	}

	tasks, amount, err := h.service.Task.FindAll(userId, columnId, pageInt, limitInt)
	if err != nil {
		return err
	}

	c.Response().Header().Set("X-Total-Count", fmt.Sprintf("%d", amount))

	return c.JSON(http.StatusOK, tasks)
}

func (h *handler) CreateTask(c echo.Context) error {
	input := new(domain.CreateTaskInput)

	if err := c.Bind(input); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	task, err := h.service.Task.Create(userId, input)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, task)
}

func (h *handler) UpdateTask(c echo.Context) error {
	input := new(domain.UpdateTaskInput)

	if err := c.Bind(input); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	taskId, err := getUIntFromParam(c, "task_id")
	if err != nil {
		return err
	}

	task, err := h.service.Task.Update(userId, taskId, input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return newResponse(c, http.StatusForbidden, err.Error())
		} else if errors.Is(err, domain.ErrRecordNotFound) {
			return newResponse(c, http.StatusNotFound, err.Error())
		}

		return err
	}

	return c.JSON(http.StatusOK, task)
}

func (h *handler) Delete(c echo.Context) error {
	taskId, err := getUIntFromParam(c, "task_id")
	if err != nil {
		return err
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	if err = h.service.Task.Delete(userId, taskId); err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return newResponse(c, http.StatusForbidden, err.Error())
		} else if errors.Is(err, domain.ErrRecordNotFound) {
			return newResponse(c, http.StatusNotFound, err.Error())
		}

		return err
	}

	return c.NoContent(http.StatusOK)
}
