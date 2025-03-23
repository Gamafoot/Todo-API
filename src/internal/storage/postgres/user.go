package postgres

import (
	"root/internal/database/models"
	"root/internal/domain"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type userStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) *userStorage {
	return &userStorage{db: db}
}

func (s *userStorage) GetById(userID uint) (*domain.User, error) {
	user := new(models.User)
	if err := s.db.First(user, "id = ?", userID).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return convert_user(user), nil
}

func (s *userStorage) GetByUsername(username string) (*domain.User, error) {
	user := new(models.User)
	if err := s.db.First(user, "username = ?", username).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return convert_user(user), nil
}

func (s *userStorage) Create(user *domain.User) error {
	if err := s.db.Create(user).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *userStorage) Delete(userId uint) error {
	if err := s.db.Delete(&domain.User{Id: userId}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func convert_user(user *models.User) *domain.User {
	return &domain.User{
		Id:        user.Id,
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}
}
