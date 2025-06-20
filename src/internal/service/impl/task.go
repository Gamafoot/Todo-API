package impl

import (
	"errors"
	"log"
	"root/internal/config"
	"root/internal/domain"
	"root/internal/storage"

	"gorm.io/gorm"
)

type taskService struct {
	config  *config.Config
	storage *storage.Storage
}

func NewTaskService(cfg *config.Config, storage *storage.Storage) *taskService {
	return &taskService{
		config:  cfg,
		storage: storage,
	}
}

func (s *taskService) List(userId, columnId uint, page, limit int) ([]*domain.Task, int, error) {
	ok, err := s.storage.Column.IsOwned(userId, columnId)
	if err != nil {
		return nil, 0, err
	}

	if !ok {
		return nil, 0, domain.ErrUserNotOwnedRecord
	}

	tasks, count, err := s.storage.Task.FindAll(columnId, page, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, domain.ErrRecordNotFound
		}

		return nil, 0, err
	}

	return tasks, count, nil
}

func (s *taskService) Create(userId uint, input *domain.CreateTaskInput) (*domain.Task, error) {
	ok, err := s.storage.Column.IsOwned(userId, input.ColumnId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	task := &domain.Task{
		ColumnId:    input.ColumnId,
		Name:        input.Name,
		Description: input.Description,
		Deadline:    input.Deadline,
	}

	taskId, err := s.storage.Task.Create(task)
	if err != nil {
		return nil, err
	}

	task, err = s.storage.Task.FindById(taskId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	err = s.storage.Heatmap.AddActivity(userId)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	return task, nil
}

func (s *taskService) Update(userId, taskId uint, input *domain.UpdateTaskInput) (*domain.Task, error) {
	ok, err := s.storage.Task.IsOwned(userId, taskId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	err = s.storage.Task.Update(&domain.Task{
		Id:          taskId,
		Name:        input.Name,
		Description: input.Description,
		Status:      input.Status,
		Archived:    input.Archived,
		Deadline:    input.Deadline,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	task, err := s.storage.Task.FindById(taskId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	if input.ColumnId > 0 && task.ColumnId != input.ColumnId {
		if input.Position <= 0 {
			input.Position = 1
		}

		err = s.storage.Task.MoveToColumn(input.ColumnId, taskId, input.Position)
		if err != nil {
			return nil, err
		}

		task.ColumnId = input.ColumnId
		task.Position = input.Position
	} else if input.Position > 0 {
		if input.Position <= 0 {
			input.Position = 1
		}

		err = s.storage.Task.MoveToPosition(task.ColumnId, taskId, input.Position)
		if err != nil {
			return nil, err
		}
		task.Position = input.Position
	}

	err = s.storage.Heatmap.AddActivity(userId)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	return task, nil
}

func (s *taskService) Delete(userId, taskId uint) error {
	ok, err := s.storage.Task.IsOwned(userId, taskId)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrUserNotOwnedRecord
	}

	if err := s.storage.Task.Delete(taskId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrRecordNotFound
		}
		return err
	}

	err = s.storage.Heatmap.RemoveActivity(userId)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	return nil
}
