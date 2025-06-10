package postgres

import (
	"root/internal/domain"
	"time"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type heatmapStorage struct {
	db *gorm.DB
}

func NewHeatmapStorage(db *gorm.DB) *heatmapStorage {
	return &heatmapStorage{db: db}
}

func (s *heatmapStorage) GetActivity(userId uint) ([]*domain.Heatmap, error) {
	heatmap := make([]*domain.Heatmap, 0)

	now := time.Now().UTC()

	err := s.db.Raw("SELECT * FROM heatmap_get_activity(?, ?)", userId, now.Year()).Scan(&heatmap).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return heatmap, nil
}

func (s *heatmapStorage) AddActivity(userId uint) error {
	err := s.db.Exec("SELECT * FROM heatmap_add_activity(?)", userId).Error
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	return nil
}

func (s *heatmapStorage) RemoveActivity(userId uint) error {
	err := s.db.Exec("SELECT * FROM heatmap_remove_activity(?)", userId).Error
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	return nil
}
