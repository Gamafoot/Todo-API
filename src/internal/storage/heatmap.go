package storage

import "root/internal/domain"

type Heatmap interface {
	GetActivity(userId uint) ([]*domain.Heatmap, error)
	AddActivity(userId uint) error
	RemoveActivity(userId uint) error
}
