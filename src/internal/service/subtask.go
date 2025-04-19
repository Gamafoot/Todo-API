package service

import (
	"root/internal/domain"
)

type SubtaskService interface {
	FindAll(userId, taskId uint, page, limit int) ([]*domain.Subtask, int, error)
	Create(userId uint, input *domain.CreateSubtaskInput) (*domain.Subtask, error)
	Update(userId, subtaskId uint, input *domain.UpdateSubtaskInput) (*domain.Subtask, error)
	Delete(userId, subtaskId uint) error
}
