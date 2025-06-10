package impl

import (
	"root/internal/config"
	"root/internal/domain"
	"root/internal/storage"
)

type heatmapService struct {
	config  *config.Config
	storage *storage.Storage
}

func NewHeatmapService(cfg *config.Config, storage *storage.Storage) *heatmapService {
	return &heatmapService{
		config:  cfg,
		storage: storage,
	}
}

func (s *heatmapService) GetActivity(userId uint) ([]*domain.Heatmap, error) {
	return s.storage.Heatmap.GetActivity(userId)
}
