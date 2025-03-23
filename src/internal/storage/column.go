package storage

import "root/internal/domain"

type ColumnStorage interface {
	FindAll(userId uint, page, limit int) ([]*domain.Column, error)
	GetAmountPages(projectId uint, page, limit int) (int, error)
	FindById(id uint) (*domain.Column, error)
	Save(column *domain.Column) error
	Delete(id uint) error
	IsOwnedUser(userId, columnId uint) (bool, error)
}
