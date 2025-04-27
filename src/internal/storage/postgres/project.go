package postgres

import (
	"math"
	"root/internal/database/model"
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

func (s *projectStorage) FindAll(userId uint, page, limit int) ([]*domain.Project, error) {
	offset := (page - 1) * limit

	projects := make([]*model.Project, 0)
	if err := s.db.Offset(offset).Limit(limit).Find(&projects, "user_id = ?", userId).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	resultProjects := make([]*domain.Project, len(projects))
	for i, project := range projects {
		resultProjects[i] = toDomainProject(project)
	}

	return resultProjects, nil
}

func (s *projectStorage) GetAmountPages(userId uint, limit int) (int, error) {
	var count int64

	if err := s.db.Model(&model.Project{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}

	amount := math.Ceil(float64(count) / float64(limit))

	return int(amount), nil
}

func (s *projectStorage) FindById(id uint) (*domain.Project, error) {
	project := new(model.Project)
	if err := s.db.Find(&project, "id = ?", id).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return toDomainProject(project), nil
}

func (s *projectStorage) Create(project *domain.Project) (uint, error) {
	modelProject := toModelProject(project)
	if err := s.db.Create(modelProject).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}
	return modelProject.Id, nil
}

func (s *projectStorage) Update(project *domain.Project) error {
	if err := s.db.Model(model.Project{}).Where("id = ?", project.Id).Updates(project).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *projectStorage) Delete(id uint) error {
	if err := s.db.Delete(&model.Project{Id: id}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *projectStorage) IsOwned(userId, projectId uint) (bool, error) {
	var isOwned bool

	err := s.db.Raw("SELECT is_owned_project(?, ?)", userId, projectId).Scan(&isOwned).Error
	if err != nil {
		return false, err
	}

	return isOwned, nil
}

func toDomainProject(project *model.Project) *domain.Project {
	return &domain.Project{
		Id:          project.Id,
		UserId:      project.UserId,
		Name:        project.Name,
		Description: project.Description,
		Deadline:    project.Deadline,
		Archived:    project.Archived,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
}

func toModelProject(project *domain.Project) *model.Project {
	return &model.Project{
		Id:          project.Id,
		UserId:      project.UserId,
		Name:        project.Name,
		Description: project.Description,
		Deadline:    project.Deadline,
		Archived:    project.Archived,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
}
