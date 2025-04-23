package postgres

import (
	"root/internal/database/model"
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
	modelSession := toModelSession(session)
	if err := s.db.Save(modelSession).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s sessionStorage) Get(userId uint, refreshToken string) (*domain.Session, error) {
	session := model.Session{}
	if err := s.db.First(&session, "user_id = ? AND refresh_token = ?", userId, refreshToken).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return toDomainSession(&session), nil
}

func toDomainSession(session *model.Session) *domain.Session {
	return &domain.Session{
		Id:           session.Id,
		UserId:       session.UserId,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
	}
}

func toModelSession(session *domain.Session) *model.Session {
	return &model.Session{
		Id:           session.Id,
		UserId:       session.UserId,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
	}
}
