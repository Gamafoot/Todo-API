package service

import (
	"errors"
	"root/internal/config"
	"root/internal/domain"
	"root/internal/storage"

	"gorm.io/gorm"
)

type ProjectService interface {
	FindAll(userId uint, page, limit int) ([]*domain.Project, int, error)
	Create(userId uint, name string) error
	Update(userId, id uint, name string) error
	Delete(userId, id uint) error
}

type projectService struct {
	config  *config.Config
	storage *storage.Storage
}

func newProjectService(cfg *config.Config, storage *storage.Storage) *projectService {
	return &projectService{
		config:  cfg,
		storage: storage,
	}
}

func (s *projectService) FindAll(userId uint, page, limit int) ([]*domain.Project, int, error) {
	projects, err := s.storage.Project.FindAll(userId, page, limit)
	if err != nil {
		return nil, 0, err
	}

	amount, err := s.storage.Project.GetAmountPages(userId, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return projects, amount, nil
}

func (s *projectService) Create(userId uint, name string) error {
	return s.storage.Project.Save(&domain.Project{
		UserId: userId,
		Name:   name,
	})
}

func (s *projectService) Update(userId, projectId uint, name string) error {
	ok, err := s.storage.Project.IsOwnedUser(userId, projectId)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrUserNotOwnedRecord
	}

	project, err := s.storage.Project.FindById(projectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrRecordNotFound
		}

		return err
	}

	project.Name = project.Name

	return s.storage.Project.Save(project)
}

func (s *projectService) Delete(userId, projectId uint) error {
	ok, err := s.storage.Project.IsOwnedUser(userId, projectId)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrUserNotOwnedRecord
	}

	return s.storage.Project.Delete(projectId)
}
