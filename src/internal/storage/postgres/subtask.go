package postgres

import (
	"errors"
	"math"
	"root/internal/database/models"
	"root/internal/domain"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type subtaskStorage struct {
	db *gorm.DB
}

func NewSubtaskStorage(db *gorm.DB) *subtaskStorage {
	return &subtaskStorage{db: db}
}

func (s *subtaskStorage) FindAll(taskId uint, page, limit int) ([]*domain.Subtask, error) {
	offset := (page - 1) * limit

	tasks := make([]*models.Subtask, 0)
	if err := s.db.Find(&tasks, "task_id = ?", taskId).Offset(offset).Limit(limit).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	resultTasks := make([]*domain.Subtask, len(tasks))
	for i, task := range tasks {
		resultTasks[i] = convertSubtask(task)
	}

	return resultTasks, nil
}

func (s *subtaskStorage) GetAmountPages(taskId uint, page, limit int) (int, error) {
	var (
		count    int64
		offset   = (page - 1) * limit
		subtasks = make([]*models.Subtask, 0)
	)

	if err := s.db.Find(&subtasks, "task_id = ?", taskId).Offset(offset).Limit(limit).Count(&count).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}

	amount := math.Ceil(float64(count) / float64(limit))

	return int(amount), nil
}

func (s *subtaskStorage) FindById(taskId uint) (*domain.Subtask, error) {
	subtask := new(models.Subtask)
	if err := s.db.Find(&subtask, "id = ?", taskId).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return convertSubtask(subtask), nil
}

func (s *subtaskStorage) Create(subtask *domain.Subtask) error {
	if err := s.db.Create(subtask).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *subtaskStorage) Update(sub *domain.Subtask) error {
	if err := s.db.Model(models.Subtask{}).Where("id = ?", sub.Id).Updates(sub).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *subtaskStorage) Delete(subtaskId uint) error {
	if err := s.db.Delete(&domain.Subtask{Id: subtaskId}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *subtaskStorage) IsOwnedUser(userId, taskId uint) (bool, error) {
	subtask := new(models.Subtask)

	err := s.db.
		Joins("JOIN tasks ON tasks.id = subtasks.task_id").
		Joins("JOIN columns ON columns.id = tasks.column_id").
		Joins("JOIN projects ON projects.id = columns.project_id").
		Joins("JOIN users ON users.id = projects.user_id").
		Where("subtasks.id = ? AND users.id = ?", taskId, userId).
		First(subtask).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func convertSubtask(task *models.Subtask) *domain.Subtask {
	return &domain.Subtask{
		Id:        task.Id,
		TaskId:    task.TaskId,
		Name:      task.Name,
		Status:    task.Status,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}
