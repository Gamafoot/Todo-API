package storage

import "root/internal/domain"

type ColumnStorage interface {
	FindAll(userId uint, page, limit int) ([]*domain.Column, error)
	GetAmountPages(projectId uint, limit int) (int, error)
	FindById(columnId uint) (*domain.Column, error)
	Create(column *domain.Column) (uint, error)
	Update(column *domain.Column) error
	Delete(columnId uint) error
	IsOwned(userId, columnId uint) (bool, error)
	MoveToPosition(projectId, columnId uint, newPosition int) error
}
