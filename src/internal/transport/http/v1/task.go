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
	api.GET("/columns/:column_id/tasks", h.FindTasks)
	api.POST("/tasks", h.CreateTask)
	api.PATCH("/tasks/:task_id", h.UpdateTask)
	api.DELETE("/tasks/:task_id", h.DeleteTask)
}

// @Summary Список задач
// @Tags task
// @Produce json
// @Security BearerAuth
// @Param column_id path int true "ID колонки"
// @Param page query int false "Номер страницы, по уполчанию 1"
// @Param limit query int false "Кол-во итоговых записей, по уполчанию 10"
// @Success 200 {array} domain.Task
// @Header 200 {integer} X-Total-Count "Общее количество задач на колонке"
// @Failure 400
// @Failure 401
// @Router /api/v1/columns/{column_id}/tasks [get]
func (h *handler) FindTasks(c echo.Context) error {
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, err.Error())
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
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return c.NoContent(http.StatusForbidden)
		}

		return err
	}

	c.Response().Header().Set("X-Total-Count", fmt.Sprintf("%d", amount))

	return c.JSON(http.StatusOK, tasks)
}

// @Summary Создать задачу
// @Tags task
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body domain.CreateTaskInput true "Данные для создания задачи"
// @Success 200 {object} domain.Task "Созданная задача"
// @Failure 400
// @Failure 401
// @Router /api/v1/tasks [post]
func (h *handler) CreateTask(c echo.Context) error {
	input := new(domain.CreateTaskInput)

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

	task, err := h.service.Task.Create(userId, input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return c.NoContent(http.StatusForbidden)
		}

		return err
	}

	return c.JSON(http.StatusOK, task)
}

// @Summary Обновить задачу
// @Tags task
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path int true "ID задачи"
// @Param body body domain.UpdateTaskInput true "Данные для обновления задачи"
// @Success 200 {object} domain.Task "Обновленная задача"
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /api/v1/tasks/{task_id} [patch]
func (h *handler) UpdateTask(c echo.Context) error {
	input := new(domain.UpdateTaskInput)

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

	taskId, err := getUIntFromParam(c, "task_id")
	if err != nil {
		return err
	}

	task, err := h.service.Task.Update(userId, taskId, input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return c.NoContent(http.StatusForbidden)
		} else if errors.Is(err, domain.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return err
	}

	return c.JSON(http.StatusOK, task)
}

// @Summary Удалить задачу
// @Tags task
// @Produce json
// @Security BearerAuth
// @Param task_id path int true "ID задачи"
// @Success 204
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /api/v1/tasks/{task_id} [delete]
func (h *handler) DeleteTask(c echo.Context) error {
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
			return c.NoContent(http.StatusForbidden)
		}

		return err
	}

	return c.NoContent(http.StatusNoContent)
}
