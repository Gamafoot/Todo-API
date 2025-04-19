package service

import (
	"root/internal/domain"
)

type ProjectService interface {
	FindAll(userId uint, page, limit int) ([]*domain.Project, int, error)
	Create(userId uint, input *domain.CreateProjectInput) (*domain.Project, error)
	Update(userId, projectId uint, input *domain.UpdateProjectInput) (*domain.Project, error)
	Delete(userId, projectId uint) error
}
