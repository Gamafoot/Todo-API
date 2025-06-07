package postgres

import (
	"root/internal/domain"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type metricsStorage struct {
	db *gorm.DB
}

func NewMetricsStorage(db *gorm.DB) *metricsStorage {
	return &metricsStorage{db: db}
}

func (s *metricsStorage) GetMetrics(projectId uint) (*domain.PreMetrics, error) {
	metrics := new(domain.PreMetrics)

	err := s.db.Raw("SELECT * FROM metrics(?)", projectId).Scan(&metrics).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return metrics, nil
}
