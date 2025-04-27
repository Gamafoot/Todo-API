package storage

import (
	"root/internal/domain"
)

type SubtaskStorage interface {
	FindAll(taskId uint, page, limit int) ([]*domain.Subtask, error)
	GetAmountPages(taskId uint, limit int) (int, error)
	FindById(subtaskId uint) (*domain.Subtask, error)
	Create(subtask *domain.Subtask) (uint, error)
	Update(subtask *domain.Subtask) error
	Delete(subtaskId uint) error
	IsOwned(userId, subtaskId uint) (bool, error)
}
