package postgres

import (
	"math"
	"root/internal/database/models"
	"root/internal/domain"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type taskStorage struct {
	db *gorm.DB
}

func NewTaskStorage(db *gorm.DB) *taskStorage {
	return &taskStorage{db: db}
}

func (s *taskStorage) GetChunk(userId uint, page, limit int) ([]*domain.Task, error) {
	offset := (page - 1) * limit

	tasks := make([]*models.Task, 0)
	if err := s.db.Find(&tasks, "user_id = ?", userId).Offset(offset).Limit(limit).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	resultTasks := make([]*domain.Task, len(tasks))
	for i, task := range tasks {
		resultTasks[i] = convert_task(task)
	}

	return resultTasks, nil
}

func (s *taskStorage) GetAmountPages(userId uint, page, limit int) (int, error) {
	var (
		count  int64
		offset = (page - 1) * limit
		tasks  = make([]*domain.Task, 0)
	)

	if err := s.db.Find(&tasks, "user_id = ?", userId).Offset(offset).Limit(limit).Count(&count).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}

	amount := math.Ceil(float64(count) / float64(limit))

	return int(amount), nil
}

func (s *taskStorage) Save(task *domain.Task) error {
	if err := s.db.Save(task).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *taskStorage) Delete(taskId, userId uint) error {
	if err := s.db.Delete(&domain.Task{}, "id = ? AND user_id = ?", taskId, userId).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func convert_task(task *models.Task) *domain.Task {
	return &domain.Task{
		Id:          task.Id,
		ProjectId:   task.ProjectId,
		Name:        task.Name,
		Description: task.Description,
		Deadline:    task.Deadline,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
