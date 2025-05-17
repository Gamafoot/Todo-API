package postgres

import (
	"root/internal/database/model"
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

func (s *columnStorage) FindAll(projectId uint, page, limit int) ([]*domain.Column, int, error) {
	offset := (page - 1) * limit

	baseQuery := s.db.Model(model.Column{})
	baseQuery = baseQuery.Where("project_id = ?", projectId)

	var count int64

	find := baseQuery.Session(&gorm.Session{})
	if err := find.Count(&count).Error; err != nil {
		return nil, 0, pkgErrors.WithStack(err)
	}

	columns := make([]*model.Column, 0)

	findQuery := baseQuery.Session(&gorm.Session{})
	findQuery = findQuery.Limit(limit).Offset(offset).Order("position")
	if err := findQuery.Find(&columns, "project_id = ?", projectId).Error; err != nil {
		return nil, 0, pkgErrors.WithStack(err)
	}

	resultColumns := make([]*domain.Column, len(columns))
	for i, column := range columns {
		resultColumns[i] = toDomainColumn(column)
	}

	return resultColumns, int(count), nil
}

func (s *columnStorage) FindById(columnId uint) (*domain.Column, error) {
	column := new(model.Column)
	if err := s.db.Find(&column, "id = ?", columnId).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return toDomainColumn(column), nil
}

func (s *columnStorage) Create(column *domain.Column) (uint, error) {
	modelColumn := toModelColumn(column)
	if err := s.db.Create(modelColumn).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}

	return modelColumn.Id, nil
}

func (s *columnStorage) Update(column *domain.Column) error {
	if err := s.db.Model(model.Column{}).Where("id = ?", column.Id).Updates(column).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *columnStorage) Delete(columnId uint) error {
	if err := s.db.Delete(&model.Column{Id: columnId}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *columnStorage) IsOwned(userId, columnId uint) (bool, error) {
	var isOwned bool

	err := s.db.Raw("SELECT is_owned_column(?, ?)", userId, columnId).Scan(&isOwned).Error
	if err != nil {
		return false, pkgErrors.WithStack(err)
	}

	return isOwned, nil
}

func (s *columnStorage) MoveToPosition(projectId, columnId uint, newPosition int) error {
	err := s.db.Exec("SELECT columns_move_to_position(?, ?, ?);", projectId, columnId, newPosition).Error
	if err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func toDomainColumn(column *model.Column) *domain.Column {
	return &domain.Column{
		Id:        column.Id,
		ProjectId: column.ProjectId,
		Name:      column.Name,
		Position:  column.Position,
	}
}

func toModelColumn(column *domain.Column) *model.Column {
	return &model.Column{
		Id:        column.Id,
		ProjectId: column.ProjectId,
		Name:      column.Name,
		Position:  column.Position,
	}
}
