package storage

import "root/internal/domain"

type TaskStorage interface {
	GetChunk(userId uint, page, limit int) ([]*domain.Task, error)
	GetAmountPages(userId uint, page, limit int) (int, error)
	Save(task *domain.Task) error
	Delete(taskId, userId uint) error
}
