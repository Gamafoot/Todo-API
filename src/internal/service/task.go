package service

import (
	"errors"
	"regexp"
	"root/internal/config"
	"root/internal/domain"
	"root/internal/storage"
	"time"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type TaskInput struct {
	UserId      uint
	Name        string
	Description string
	Status      string
	Deadline    string
}

type TaskService interface {
	GetChunk(userId uint, page, limit int) ([]*domain.Task, int, error)
	Save(inp *TaskInput) error
	Delete(userId, taskId uint) error
}

type taskService struct {
	config  *config.Config
	storage *storage.Storage
}

func newTaskService(cfg *config.Config, storage *storage.Storage) *taskService {
	return &taskService{
		config:  cfg,
		storage: storage,
	}
}

func (s *taskService) GetChunk(userId uint, page, limit int) ([]*domain.Task, int, error) {
	tasks, err := s.storage.Task.GetChunk(userId, page, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, domain.ErrTaskNotFound
		}

		return nil, 0, err
	}

	amount, err := s.storage.Task.GetAmountPages(userId, page, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, domain.ErrTaskNotFound
		}

		return nil, 0, err
	}

	return tasks, amount, nil
}

func (s *taskService) Save(inp *TaskInput) error {
	var (
		deadline time.Time
		err      error
	)

	if len(inp.Deadline) > 0 {
		deadline, err = parseDeadline(inp.Deadline)
		if err != nil {
			return err
		}
	}

	err = s.storage.Task.Save(&domain.Task{
		ProjectId:   inp.UserId, // fail!!!
		Name:        inp.Name,
		Description: inp.Description,
		Status:      inp.Status,
		Deadline:    &deadline,
	})
	if err != nil {
		return err
	}

	return nil
}

func parseDeadline(deadline string) (time.Time, error) {
	var result time.Time

	ok, err := regexp.MatchString(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}`, deadline)
	if err != nil {
		return result, pkgErrors.WithStack(err)
	}

	if ok {
		result, err = time.Parse("2006-01-02 15:04", deadline)
		if err != nil {
			return result, pkgErrors.WithStack(err)
		}

		return result, nil
	}

	ok, err = regexp.MatchString(`\d{4}-\d{2}-\d{2}`, deadline)
	if err != nil {
		return result, pkgErrors.WithStack(err)
	}

	if ok {
		result, err = time.Parse("2006-01-02", deadline)
		if err != nil {
			return result, pkgErrors.WithStack(err)
		}

		return result, nil
	}

	return result, domain.ErrInvalidDeadlineFormat
}

func (s *taskService) Delete(taskId, userId uint) error {
	if err := s.storage.Task.Delete(taskId, userId); err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrTaskNotFound
		}

		return err
	}

	return nil
}
