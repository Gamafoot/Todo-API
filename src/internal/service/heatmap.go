package service

import "root/internal/domain"

type Heatmap interface {
	GetActivity(userId uint) ([]*domain.Heatmap, error)
}
