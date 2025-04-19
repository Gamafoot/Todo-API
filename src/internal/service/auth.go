package service

import (
	"root/internal/domain"
)

//go:generate mockgen -source=auth.go -destination=mocks/mock_auth.go

type AuthService interface {
	Login(input *domain.LoginInput) (*domain.Tokens, error)
	Register(input *domain.RegisterInput) error
	RefreshToken(userId uint, refreshToken string) (*domain.Tokens, error)
}
