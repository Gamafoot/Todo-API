package service

import (
	"root/internal/domain"
)

type TaskService interface {
	FindAll(userId, columnId uint, page, limit int) ([]*domain.Task, int, error)
	Create(userId uint, input *domain.CreateTaskInput) (*domain.Task, error)
	Update(userId, taskId uint, input *domain.UpdateTaskInput) (*domain.Task, error)
	Delete(userId, taskId uint) error
}
