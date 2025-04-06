package storage

import (
	"root/internal/domain"
)

type TaskStorage interface {
	FindAll(columnId uint, page, limit int) ([]*domain.Task, error)
	GetAmountPages(userId uint, page, limit int) (int, error)
	FindById(taskId uint) (*domain.Task, error)
	Create(task *domain.Task) error
	Update(task *domain.Task) error
	Delete(taskId uint) error
	IsOwnedUser(userId, taskId uint) (bool, error)
}
