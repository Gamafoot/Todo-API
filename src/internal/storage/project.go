package storage

import "root/internal/domain"

type ProjectStorage interface {
	GetById(id uint) (*domain.Project, error)
	Save(project *domain.Project) error
	Delete(id uint) error
}
