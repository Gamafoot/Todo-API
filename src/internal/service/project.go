package service

import (
	"root/internal/domain"
)

type ProjectService interface {
	List(userId uint, options *domain.SearchProjectOptions, page, limit int) ([]*domain.Project, int, error)
	Detail(userId, projectId uint) (*domain.Project, error)
	Create(userId uint, input *domain.CreateProjectInput) (*domain.Project, error)
	Update(userId, projectId uint, input *domain.UpdateProjectInput) (*domain.Project, error)
	Delete(userId, projectId uint) error

	ProjectStatsService
	ProjectMetrics
}

type ProjectStatsService interface {
	GetStats(userId, projectId uint) (*domain.ProjectStats, error)
}

type ProjectMetrics interface {
	GetMetrics(userId, projectId uint) (*domain.Metrics, error)
}
