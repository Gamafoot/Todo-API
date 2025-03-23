package storage

import "root/internal/domain"

type ProjectStorage interface {
	FindAll(userId uint, page, limit int) ([]*domain.Project, error)
	GetAmountPages(userId uint, page, limit int) (int, error)
	FindById(projectId uint) (*domain.Project, error)
	Save(project *domain.Project) error
	Delete(id uint) error
	IsOwnedUser(userId, projectId uint) (bool, error)
}
