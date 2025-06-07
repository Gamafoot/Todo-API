package storage

import "root/internal/domain"

type ProjectStorage interface {
	FindAll(userId uint, options *domain.SearchProjectOptions, page, limit int) ([]*domain.Project, int, error)
	FindById(projectId uint) (*domain.Project, error)
	Create(project *domain.Project) (uint, error)
	Update(project *domain.Project) error
	Delete(projectId uint) error
	IsOwned(userId, projectId uint) (bool, error)

	ProjectStatsStorage
}

type ProjectStatsStorage interface {
	GetStats(projectId uint) (*domain.ProjectStats, error)
	GetMetrics(projectId uint) (*domain.PreProjectMetrics, error)
	GetProgress(projectId uint) ([]*domain.ProjectProgress, error)
}
