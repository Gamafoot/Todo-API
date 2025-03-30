package postgres

import (
	"errors"
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

func (s *taskStorage) FindAll(columnId uint, page, limit int) ([]*domain.Task, error) {
	offset := (page - 1) * limit

	tasks := make([]*models.Task, 0)
	if err := s.db.Find(&tasks, "column_id = ?", columnId).Offset(offset).Limit(limit).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	resultTasks := make([]*domain.Task, len(tasks))
	for i, task := range tasks {
		resultTasks[i] = convertTask(task)
	}

	return resultTasks, nil
}

func (s *taskStorage) GetAmountPages(columnId uint, page, limit int) (int, error) {
	var (
		count  int64
		offset = (page - 1) * limit
		tasks  = make([]*models.Task, 0)
	)

	if err := s.db.Find(&tasks, "column_id = ?", columnId).Offset(offset).Limit(limit).Count(&count).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}

	amount := math.Ceil(float64(count) / float64(limit))

	return int(amount), nil
}

func (s *taskStorage) FindById(taskId uint) (*domain.Task, error) {
	task := new(models.Task)
	if err := s.db.Find(&task, "id = ?", taskId).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return convertTask(task), nil
}

func (s *taskStorage) Create(task *domain.Task) error {
	if err := s.db.Create(task).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *taskStorage) Update(task *domain.Task) error {
	if err := s.db.Model(models.Task{}).Where("id = ?", task.Id).Updates(task).Error; err != nil {
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

func (s *taskStorage) IsOwnedUser(userId, columnId uint) (bool, error) {
	column := new(models.Column)

	err := s.db.
		Joins("JOIN columns ON columns.id = tasks.column_id").
		Joins("JOIN projects ON projects.id = columns.project_id").
		Joins("JOIN users ON users.id = projects.user_id").
		Where("column.id = ? AND users.id = ?", columnId, userId).
		First(column).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func convertTask(task *models.Task) *domain.Task {
	return &domain.Task{
		Id:          task.Id,
		ColumnId:    task.ColumnId,
		Name:        task.Name,
		Description: task.Description,
		Deadline:    task.Deadline,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
