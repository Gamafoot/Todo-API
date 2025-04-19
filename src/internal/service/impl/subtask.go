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

func (s *subtaskService) FindAll(userId, taskId uint, page, limit int) ([]*domain.Subtask, int, error) {
	ok, err := s.storage.Task.IsOwnedUser(userId, taskId)
	if err != nil {
		return nil, 0, err
	}

	if !ok {
		return nil, 0, domain.ErrUserNotOwnedRecord
	}

	tasks, err := s.storage.Subtask.FindAll(taskId, page, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, domain.ErrRecordNotFound
		}

		return nil, 0, err
	}

	amount, err := s.storage.Subtask.GetAmountPages(taskId, page, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, domain.ErrRecordNotFound
		}

		return nil, 0, err
	}

	return tasks, int(amount), nil
}

func (s *subtaskService) Create(userId uint, input *domain.CreateSubtaskInput) (*domain.Subtask, error) {
	ok, err := s.storage.Task.IsOwnedUser(userId, input.TaskId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	task := &domain.Subtask{
		TaskId:   input.TaskId,
		Name:     input.Name,
		Status:   input.Status,
		Deadline: input.Deadline,
	}

	err = s.storage.Subtask.Create(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *subtaskService) Update(userId, subtaskId uint, input *domain.UpdateSubtaskInput) (*domain.Subtask, error) {
	ok, err := s.storage.Subtask.IsOwnedUser(userId, subtaskId)
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
		Deadline: input.Deadline,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}

		return nil, err
	}

	task, err := s.storage.Subtask.FindById(subtaskId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}

		return nil, err
	}

	return task, nil
}

func (s *subtaskService) Delete(userId, subtaskId uint) error {
	ok, err := s.storage.Subtask.IsOwnedUser(userId, subtaskId)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrUserNotOwnedRecord
	}

	return s.storage.Subtask.Delete(subtaskId)
}
