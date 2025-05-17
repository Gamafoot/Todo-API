package storage

import (
	"root/internal/domain"
)

type TaskStorage interface {
	FindAll(columnId uint, page, limit int) ([]*domain.Task, int, error)
	FindById(taskId uint) (*domain.Task, error)
	Create(task *domain.Task) (uint, error)
	Update(task *domain.Task) error
	Delete(taskId uint) error
	IsOwned(userId, taskId uint) (bool, error)
	MoveToColumn(columnId, taskId uint, newPosition int) error
	MoveToPosition(columnId, taskId uint, newPosition int) error
}
