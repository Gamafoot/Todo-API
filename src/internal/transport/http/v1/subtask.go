package v1

import (
	"errors"
	"fmt"
	"net/http"
	"root/internal/domain"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) initSubTaskRoutes(api *echo.Group) {
	api.GET("/tasks/:task_id/subtasks", h.FindSubtasks)
	api.POST("/subtasks", h.CreateSubtask)
	api.PATCH("/subtasks/:subtask_id", h.UpdateSubtask)
	api.DELETE("/subtasks/:subtask_id", h.DeleteSubtask)
}

// @Summary Список подзадач
// @Tags subtask
// @Produce json
// @Security BearerAuth
// @Param task_id path int true "ID задачи"
// @Param page query int false "Номер страницы, по уполчанию 1"
// @Param limit query int false "Кол-во итоговых записей, по уполчанию 10"
// @Success 200 {array} domain.Subtask
// @Header 200 {integer} X-Total-Count "Общее количество подзадач на колонке"
// @Failure 400
// @Router /api/v1/tasks/{task_id}/subtasks [get]
func (h *handler) FindSubtasks(c echo.Context) error {
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

	taskId, err := getUIntFromParam(c, "task_id")
	if err != nil {
		return err
	}

	tasks, amount, err := h.service.Subtask.FindAll(userId, taskId, pageInt, limitInt)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return c.NoContent(http.StatusForbidden)
		}

		return err
	}

	c.Response().Header().Set("X-Total-Count", fmt.Sprintf("%d", amount))

	return c.JSON(http.StatusOK, tasks)
}

// @Summary Создать подзадачу
// @Tags subtask
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body domain.CreateSubtaskInput true "Данные для создания подзадачи"
// @Success 200 {object} domain.Subtask "Созданная подзадача"
// @Failure 400
// @Router /api/v1/subtasks [post]
func (h *handler) CreateSubtask(c echo.Context) error {
	input := new(domain.CreateSubtaskInput)

	if err := c.Bind(input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	subtask, err := h.service.Subtask.Create(userId, input)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return c.NoContent(http.StatusForbidden)
		}

		return err
	}

	return c.JSON(http.StatusOK, subtask)
}

// @Summary Обновить подзадачу
// @Tags subtask
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param subtask_id path int true "ID подзадачи"
// @Param body body domain.UpdateTaskInput true "Данные для обновления подзадачи"
// @Success 200 {object} domain.Subtask "Обновленная подзадача"
// @Failure 400
// @Failure 403
// @Failure 404
// @Router /api/v1/subtasks/{subtask_id} [patch]
func (h *handler) UpdateSubtask(c echo.Context) error {
	input := new(domain.UpdateSubtaskInput)

	if err := c.Bind(input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	subtaskId, err := getUIntFromParam(c, "subtask_id")
	if err != nil {
		return err
	}

	task, err := h.service.Subtask.Update(userId, subtaskId, input)
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

// @Summary Удалить подзадачу
// @Tags subtask
// @Produce json
// @Security BearerAuth
// @Param subtask_id path int true "ID подзадачи"
// @Success 204
// @Failure 403
// @Failure 404
// @Router /api/v1/subtasks/{subtask_id} [delete]
func (h *handler) DeleteSubtask(c echo.Context) error {
	subtaskId, err := getUIntFromParam(c, "subtask_id")
	if err != nil {
		return err
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	if err = h.service.Subtask.Delete(userId, subtaskId); err != nil {
		if errors.Is(err, domain.ErrUserNotOwnedRecord) {
			return c.NoContent(http.StatusForbidden)
		}

		return err
	}

	return c.NoContent(http.StatusNoContent)
}
