package postgres

import (
	"root/internal/database/models"
	"root/internal/domain"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type projectStorage struct {
	db *gorm.DB
}

func NewProjectStorage(db *gorm.DB) *projectStorage {
	return &projectStorage{db: db}
}

func (s *projectStorage) GetById(id uint) (*domain.Project, error) {
	project := new(models.Project)
	if err := s.db.First(project, "id = ?", id).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return convert_project(project), nil
}

func (s *projectStorage) Save(project *domain.Project) error {
	d := models.Project{
		Id: project.Id,
		Name: project.Name,
	}
	if err := s.db.Save(project).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *projectStorage) Delete(id uint) error {
	if err := s.db.Delete(&domain.Project{Id: id}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func convert_project(project *models.Project) *domain.Project {
	return &domain.Project{
		Id:   project.Id,
		Name: project.Name,
	}
}
