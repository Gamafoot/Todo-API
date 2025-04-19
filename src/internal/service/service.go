package service

import (
	"root/internal/config"
	"root/internal/service/impl"
	"root/internal/storage"
	"root/pkg/jwt"
)

type Service struct {
	Auth    AuthService
	Project ProjectService
	Column  ColumnService
	Task    TaskService
	Subtask SubtaskService
}

func NewService(cfg *config.Config, storage *storage.Storage, tokenManager jwt.TokenManager) *Service {
	return &Service{
		Auth:    impl.NewAuthService(cfg, storage, tokenManager),
		Project: impl.NewProjectService(cfg, storage),
		Column:  impl.NewColumnService(cfg, storage),
		Task:    impl.NewTaskService(cfg, storage),
		Subtask: impl.NewSubtaskService(cfg, storage),
	}
}
