package storage

import (
	"root/internal/domain"
)

type MetricsStorage interface {
	GetMetrics(projectId uint) (*domain.PreMetrics, error)
}
