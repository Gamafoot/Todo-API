package service

import (
	"root/internal/config"
	"root/internal/storage"
	"root/pkg/jwt"
)

type Service struct {
	Auth AuthService
	Task TaskService
}

func NewService(cfg *config.Config, storage *storage.Storage, tokenManager jwt.TokenManager) *Service {
	return &Service{
		Auth: newAuthService(cfg, storage, tokenManager),
		Task: newTaskService(cfg, storage),
	}
}
