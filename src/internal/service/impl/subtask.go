package impl

import (
	"errors"
	"root/internal/config"
	"root/internal/domain"
	"root/internal/storage"

	"gorm.io/gorm"
)

type subtaskService struct {
	config  *config.Config
	storage *storage.Storage
}

func NewSubtaskService(cfg *config.Config, storage *storage.Storage) *subtaskService {
	return &subtaskService{
		config:  cfg,
		storage: storage,
	}
}

func (s *subtaskService) List(userId, taskId uint, page, limit int) ([]*domain.Subtask, int, error) {
	ok, err := s.storage.Task.IsOwned(userId, taskId)
	if err != nil {
		return nil, 0, err
	}

	if !ok {
		return nil, 0, domain.ErrUserNotOwnedRecord
	}

	tasks, count, err := s.storage.Subtask.FindAll(taskId, page, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, domain.ErrRecordNotFound
		}

		return nil, 0, err
	}

	return tasks, count, nil
}

func (s *subtaskService) Create(userId uint, input *domain.CreateSubtaskInput) (*domain.Subtask, error) {
	ok, err := s.storage.Task.IsOwned(userId, input.TaskId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	subtask := &domain.Subtask{
		TaskId: input.TaskId,
		Name:   input.Name,
	}

	subtaskId, err := s.storage.Subtask.Create(subtask)
	if err != nil {
		return nil, err
	}

	subtask, err = s.storage.Subtask.FindById(subtaskId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	return subtask, nil
}

func (s *subtaskService) Update(userId, subtaskId uint, input *domain.UpdateSubtaskInput) (*domain.Subtask, error) {
	ok, err := s.storage.Subtask.IsOwned(userId, subtaskId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	err = s.storage.Subtask.Update(&domain.Subtask{
		Id:       subtaskId,
		Name:     input.Name,
		Status:   input.Status,
		Archived: input.Archived,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	subtask, err := s.storage.Subtask.FindById(subtaskId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}

		return nil, err
	}

	if input.Position > 0 {
		err = s.storage.Subtask.MoveToPosition(subtask.TaskId, subtaskId, input.Position)
		if err != nil {
			return nil, err
		}
		subtask.Position = input.Position
	}

	return subtask, nil
}

func (s *subtaskService) Delete(userId, subtaskId uint) error {
	ok, err := s.storage.Subtask.IsOwned(userId, subtaskId)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrUserNotOwnedRecord
	}

	if err := s.storage.Subtask.Delete(subtaskId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrRecordNotFound
		}
		return err
	}

	return nil
}
