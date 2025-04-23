package postgres

import (
	"math"
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

func (s *projectStorage) FindAll(userId uint, page, limit int) ([]*domain.Project, error) {
	offset := (page - 1) * limit

	projects := make([]*models.Project, 0)
	if err := s.db.Find(&projects, "user_id = ?", userId).Offset(offset).Limit(limit).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	resultProjects := make([]*domain.Project, len(projects))
	for i, project := range projects {
		resultProjects[i] = convertProject(project)
	}

	return resultProjects, nil
}

func (s *projectStorage) GetAmountPages(userId uint, page, limit int) (int, error) {
	var (
		count    int64
		offset   = (page - 1) * limit
		projects = make([]*models.Project, 0)
	)

	if err := s.db.Find(&projects, "user_id = ?", userId).Offset(offset).Limit(limit).Count(&count).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}

	amount := math.Ceil(float64(count) / float64(limit))

	return int(amount), nil
}

func (s *projectStorage) FindById(id uint) (*domain.Project, error) {
	project := new(models.Project)
	if err := s.db.Find(&project, "id = ?", id).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return convertProject(project), nil
}

func (s *projectStorage) Create(project *domain.Project) error {
	if err := s.db.Create(project).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *projectStorage) Update(project *domain.Project) error {
	if err := s.db.Model(models.Project{}).Where("id = ?", project.Id).Updates(project).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *projectStorage) Delete(id uint) error {
	if err := s.db.Delete(&models.Project{Id: id}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *projectStorage) IsOwnedUser(userId, projectId uint) (bool, error) {
	project := new(models.Project)

	err := s.db.
		Joins("JOIN users ON users.id = projects.user_id").
		Where("users.id = ? AND projects.id = ?", userId, projectId).
		First(project).
		Error

	if err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func convertProject(project *models.Project) *domain.Project {
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
