package postgres

import (
	"math"
	"root/internal/database/model"
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

	tasks := make([]*model.Task, 0)
	if err := s.db.Offset(offset).Limit(limit).Order("position").Find(&tasks, "column_id = ?", columnId).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	resultTasks := make([]*domain.Task, len(tasks))
	for i, task := range tasks {
		resultTasks[i] = toDomainTask(task)
	}

	return resultTasks, nil
}

func (s *taskStorage) GetAmountPages(columnId uint, limit int) (int, error) {
	var count int64

	if err := s.db.Model(&model.Task{}).Where("column_id = ?", columnId).Count(&count).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}

	amount := math.Ceil(float64(count) / float64(limit))

	return int(amount), nil
}

func (s *taskStorage) FindById(taskId uint) (*domain.Task, error) {
	task := new(model.Task)
	if err := s.db.Find(&task, "id = ?", taskId).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return toDomainTask(task), nil
}

func (s *taskStorage) Create(task *domain.Task) (uint, error) {
	modelTask := toModelTask(task)
	if err := s.db.Create(modelTask).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}
	return modelTask.Id, nil
}

func (s *taskStorage) Update(task *domain.Task) error {
	if err := s.db.Model(model.Task{}).Where("id = ?", task.Id).Updates(task).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *taskStorage) Delete(taskId uint) error {
	if err := s.db.Delete(&domain.Task{Id: taskId}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *taskStorage) IsOwned(userId, taskId uint) (bool, error) {
	var isOwned bool

	err := s.db.Raw("SELECT is_owned_task(?, ?)", userId, taskId).Scan(&isOwned).Error
	if err != nil {
		return false, pkgErrors.WithStack(err)
	}

	return isOwned, nil
}

func (s *taskStorage) MoveToColumn(columnId, taskId uint, newPosition int) error {
	err := s.db.Exec("SELECT tasks_move_to_column(?, ?, ?);", columnId, taskId, newPosition).Error
	if err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *taskStorage) MoveToPosition(columnId, taskId uint, newPosition int) error {
	err := s.db.Exec("SELECT tasks_move_to_position(?, ?, ?);", columnId, taskId, newPosition).Error
	if err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func toDomainTask(task *model.Task) *domain.Task {
	return &domain.Task{
		Id:          task.Id,
		ColumnId:    task.ColumnId,
		Name:        task.Name,
		Status:      task.Status,
		Archived:    task.Archived,
		Position:    task.Position,
		Description: task.Description,
		Deadline:    task.Deadline,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func toModelTask(task *domain.Task) *model.Task {
	return &model.Task{
		Id:          task.Id,
		ColumnId:    task.ColumnId,
		Name:        task.Name,
		Status:      task.Status,
		Archived:    task.Archived,
		Position:    task.Position,
		Description: task.Description,
		Deadline:    task.Deadline,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
