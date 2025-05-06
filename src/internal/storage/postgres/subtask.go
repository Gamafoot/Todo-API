package postgres

import (
	"math"
	"root/internal/database/model"
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

	tasks := make([]*model.Subtask, 0)
	if err := s.db.Offset(offset).Limit(limit).Order("position").Find(&tasks, "task_id = ?", taskId).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	resultTasks := make([]*domain.Subtask, len(tasks))
	for i, task := range tasks {
		resultTasks[i] = toDomainSubtask(task)
	}

	return resultTasks, nil
}

func (s *subtaskStorage) GetAmountPages(taskId uint, limit int) (int, error) {
	var count int64

	if err := s.db.Model(&model.Subtask{}).Where("task_id = ?", taskId).Count(&count).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}

	amount := math.Ceil(float64(count) / float64(limit))

	return int(amount), nil
}

func (s *subtaskStorage) FindById(taskId uint) (*domain.Subtask, error) {
	subtask := new(model.Subtask)
	if err := s.db.Find(&subtask, "id = ?", taskId).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return toDomainSubtask(subtask), nil
}

func (s *subtaskStorage) Create(subtask *domain.Subtask) (uint, error) {
	modelSubtask := toModelSubtask(subtask)
	if err := s.db.Create(modelSubtask).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}
	return modelSubtask.Id, nil
}

func (s *subtaskStorage) Update(sub *domain.Subtask) error {
	if err := s.db.Model(model.Subtask{}).Where("id = ?", sub.Id).Updates(sub).Error; err != nil {
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

func (s *subtaskStorage) IsOwned(userId, subtaskId uint) (bool, error) {
	var isOwned bool

	err := s.db.Raw("SELECT is_owned_subtask(?, ?)", userId, subtaskId).Scan(&isOwned).Error
	if err != nil {
		return false, pkgErrors.WithStack(err)
	}

	return isOwned, nil
}

func (s *subtaskStorage) MoveToPosition(taskId, subtaskId uint, newPosition int) error {
	err := s.db.Exec("SELECT subtasks_move_to_position(?, ?, ?);", taskId, subtaskId, newPosition).Error
	if err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func toDomainSubtask(subtask *model.Subtask) *domain.Subtask {
	return &domain.Subtask{
		Id:        subtask.Id,
		TaskId:    subtask.TaskId,
		Name:      subtask.Name,
		Status:    subtask.Status,
		Archived:  subtask.Archived,
		Position:  subtask.Position,
		CreatedAt: subtask.CreatedAt,
		UpdatedAt: subtask.UpdatedAt,
	}
}

func toModelSubtask(subtask *domain.Subtask) *model.Subtask {
	return &model.Subtask{
		Id:        subtask.Id,
		TaskId:    subtask.TaskId,
		Name:      subtask.Name,
		Status:    subtask.Status,
		Archived:  subtask.Archived,
		Position:  subtask.Position,
		CreatedAt: subtask.CreatedAt,
		UpdatedAt: subtask.UpdatedAt,
	}
}
