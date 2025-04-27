package impl

import (
	"errors"
	"root/internal/config"
	"root/internal/domain"
	"root/internal/storage"

	"gorm.io/gorm"
)

type columnService struct {
	config  *config.Config
	storage *storage.Storage
}

func NewColumnService(cfg *config.Config, storage *storage.Storage) *columnService {
	return &columnService{
		config:  cfg,
		storage: storage,
	}
}

func (s *columnService) List(userId, projectId uint, page, limit int) ([]*domain.Column, int, error) {
	ok, err := s.storage.Project.IsOwned(userId, projectId)
	if err != nil {
		return nil, 0, err
	}

	if !ok {
		return nil, 0, domain.ErrUserNotOwnedRecord
	}

	columns, err := s.storage.Column.FindAll(projectId, page, limit)
	if err != nil {
		return nil, 0, err
	}

	amount, err := s.storage.Task.GetAmountPages(userId, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, domain.ErrRecordNotFound
		}
		return nil, 0, err
	}

	return columns, amount, nil
}

func (s *columnService) Create(userId uint, input *domain.CreateColumnInput) (*domain.Column, error) {
	ok, err := s.storage.Project.IsOwned(userId, input.ProjectId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	column := &domain.Column{
		ProjectId: input.ProjectId,
		Name:      input.Name,
	}

	columnId, err := s.storage.Column.Create(column)
	if err != nil {
		return nil, err
	}

	column, err = s.storage.Column.FindById(columnId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	return column, nil
}

func (s *columnService) Update(userId, columnId uint, input *domain.UpdateColumnInput) (*domain.Column, error) {
	ok, err := s.storage.Column.IsOwned(userId, columnId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	err = s.storage.Column.Update(&domain.Column{
		Id:   columnId,
		Name: input.Name,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	column, err := s.storage.Column.FindById(columnId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	return column, nil
}

func (s *columnService) Delete(userId, columnId uint) error {
	ok, err := s.storage.Column.IsOwned(userId, columnId)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrUserNotOwnedRecord
	}

	if err = s.storage.Column.Delete(columnId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrRecordNotFound
		}
		return err
	}

	return nil
}
