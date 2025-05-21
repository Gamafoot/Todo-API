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

func (s *projectService) List(userId uint, options *domain.SearchProjectOptions, page, limit int) ([]*domain.Project, int, error) {
	return s.storage.Project.FindAll(userId, options, page, limit)
}

func (s *projectService) Detail(userId, projectId uint) (*domain.Project, error) {
	ok, err := s.storage.Project.IsOwned(userId, projectId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	project, err := s.storage.Project.FindById(projectId)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) Create(userId uint, input *domain.CreateProjectInput) (*domain.Project, error) {
	project := &domain.Project{
		UserId:      userId,
		Name:        input.Name,
		Description: input.Description,
		Deadline:    input.Deadline,
	}

	projectId, err := s.storage.Project.Create(project)
	if err != nil {
		return nil, err
	}

	project, err = s.storage.Project.FindById(projectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	return project, nil
}

func (s *projectService) Update(userId, projectId uint, input *domain.UpdateProjectInput) (*domain.Project, error) {
	ok, err := s.storage.Project.IsOwned(userId, projectId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	err = s.storage.Project.Update(&domain.Project{
		Id:          projectId,
		Name:        input.Name,
		Description: input.Description,
		Archived:    input.Archived,
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
	ok, err := s.storage.Project.IsOwned(userId, projectId)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrUserNotOwnedRecord
	}

	if err := s.storage.Project.Delete(projectId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrRecordNotFound
		}
		return err
	}

	return nil
}

func (s *projectService) GetStats(userId, projectId uint) (*domain.ProjectStats, error) {
	ok, err := s.storage.Project.IsOwned(userId, projectId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	return s.storage.Project.GetStats(projectId)
}
