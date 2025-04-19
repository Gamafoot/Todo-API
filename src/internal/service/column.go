package service

import (
	"root/internal/domain"
)

type ColumnService interface {
	FindAll(userId, projectId uint, page, limit int) ([]*domain.Column, int, error)
	Create(userId uint, input *domain.CreateColumnInput) (*domain.Column, error)
	Update(userId, columnId uint, input *domain.UpdateColumnInput) (*domain.Column, error)
	Delete(userId, columnId uint) error
}
