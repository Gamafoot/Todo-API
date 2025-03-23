package service

import (
	"root/internal/config"
	"root/internal/storage"
	"root/pkg/jwt"
)

type Service struct {
	Auth AuthService
	Task TaskService
	Project ProjectService
	Column ColumnService
}

func NewService(cfg *config.Config, storage *storage.Storage, tokenManager jwt.TokenManager) *Service {
	return &Service{
		Auth: newAuthService(cfg, storage, tokenManager),
		Task: newTaskService(cfg, storage),
		Project: newProjectService(cfg, storage),
		Column: newColumnService(cfg, storage),
	}
}
