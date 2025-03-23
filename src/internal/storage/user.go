package storage

import "root/internal/domain"

type UserStorage interface {
	GetById(id uint) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	Create(user *domain.User) error
	Delete(id uint) error
}
