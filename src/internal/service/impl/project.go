package impl

import (
	"errors"
	"root/internal/config"
	"root/internal/domain"
	"root/internal/storage"

	"gorm.io/gorm"
)

type projectService struct {
	config  *config.Config
	storage *storage.Storage
}

func NewProjectService(cfg *config.Config, storage *storage.Storage) *projectService {
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

func (s *projectService) Create(userId uint, input *domain.CreateProjectInput) (*domain.Project, error) {
	project := &domain.Project{
		UserId: userId,
		Name:   input.Name,
	}

	err := s.storage.Project.Create(project)
	if err != nil {
		return project, nil
	}

	return project, nil
}

func (s *projectService) Update(userId, projectId uint, input *domain.UpdateProjectInput) (*domain.Project, error) {
	ok, err := s.storage.Project.IsOwnedUser(userId, projectId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	err = s.storage.Project.Update(&domain.Project{
		Id:   projectId,
		Name: input.Name,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}

		return nil, err
	}

	project, err := s.storage.Project.FindById(projectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	return project, nil
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
