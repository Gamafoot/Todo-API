package postgres

import (
	"root/internal/database/models"
	"root/internal/domain"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type sessionStorage struct {
	db *gorm.DB
}

func NewSessionStorage(db *gorm.DB) *sessionStorage {
	return &sessionStorage{db: db}
}

func (s sessionStorage) Set(session *domain.Session) error {
	if err := s.db.Save(session).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s sessionStorage) Get(userId uint, refreshToken string) (domain.Session, error) {
	session := models.Session{}
	if err := s.db.First(&session, "user_id = ? AND refresh_token = ?", userId, refreshToken).Error; err != nil {
		return convertSession(&session), pkgErrors.WithStack(err)
	}

	return convertSession(&session), nil
}

func convertSession(session *models.Session) domain.Session {
	return domain.Session{
		Id:           session.Id,
		UserId:       session.UserId,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
	}
}
