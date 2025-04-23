package storage

import "root/internal/domain"

type SessionStorage interface {
	Set(session *domain.Session) error
	Get(userId uint, refreshToken string) (*domain.Session, error)
}
