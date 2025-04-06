package postgres

import (
	"errors"
	"math"
	"root/internal/database/models"
	"root/internal/domain"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type columnStorage struct {
	db *gorm.DB
}

func NewColumnStorage(db *gorm.DB) *columnStorage {
	return &columnStorage{db: db}
}

func (s *columnStorage) FindAll(projectId uint, page, limit int) ([]*domain.Column, error) {
	offset := (page - 1) * limit

	columns := make([]*models.Column, 0)
	if err := s.db.Find(&columns, "project_id = ?", projectId).Offset(offset).Limit(limit).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	resultColumns := make([]*domain.Column, len(columns))
	for i, column := range columns {
		resultColumns[i] = convertColumn(column)
	}

	return resultColumns, nil
}

func (s *columnStorage) GetAmountPages(columnId uint, page, limit int) (int, error) {
	var (
		count  int64
		offset = (page - 1) * limit
		tasks  = make([]*models.Column, 0)
	)

	if err := s.db.Find(&tasks, "column_id = ?", columnId).Offset(offset).Limit(limit).Count(&count).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}

	amount := math.Ceil(float64(count) / float64(limit))

	return int(amount), nil
}

func (s *columnStorage) FindById(columnId uint) (*domain.Column, error) {
	column := new(models.Column)
	if err := s.db.Find(&column, "id = ?", columnId).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return convertColumn(column), nil
}

func (s *columnStorage) Create(column *domain.Column) error {
	if err := s.db.Create(column).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *columnStorage) Update(column *domain.Column) error {
	if err := s.db.Model(models.Column{}).Where("id = ?", column.Id).Updates(column).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *columnStorage) Delete(columnId uint) error {
	if err := s.db.Delete(&models.Column{Id: columnId}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *columnStorage) IsOwnedUser(userId, columnId uint) (bool, error) {
	column := new(models.Column)

	err := s.db.
		Joins("JOIN projects ON projects.id = columns.project_id").
		Joins("JOIN users ON users.id = projects.user_id").
		Where("columns.id = ? AND users.id = ?", columnId, userId).
		First(column).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func convertColumn(column *models.Column) *domain.Column {
	return &domain.Column{
		Id:        column.Id,
		ProjectId: column.ProjectId,
		Name:      column.Name,
	}
}
