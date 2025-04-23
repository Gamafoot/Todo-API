package storage

import "root/internal/domain"

type ProjectStorage interface {
	FindAll(userId uint, page, limit int) ([]*domain.Project, error)
	GetAmountPages(userId uint, limit int) (int, error)
	FindById(projectId uint) (*domain.Project, error)
	Create(project *domain.Project) (uint, error)
	Update(project *domain.Project) error
	Delete(projectId uint) error
	IsOwnedUser(userId, projectId uint) (bool, error)
}
