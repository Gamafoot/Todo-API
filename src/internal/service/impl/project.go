package impl

import (
	"errors"
	"log"
	"math"
	"root/internal/config"
	"root/internal/domain"
	"root/internal/storage"
	"time"

	"gorm.io/gorm"
)

type projectService struct {
	config  *config.Config
	storage *storage.Storage
}

func NewProjectService(cfg *config.Config, storage *storage.Storage) *projectService {
	return &projectService{
		config:  cfg,
		storage: storage,
	}
}

func (s *projectService) List(userId uint, options *domain.SearchProjectOptions, page, limit int) ([]*domain.Project, int, error) {
	return s.storage.Project.FindAll(userId, options, page, limit)
}

func (s *projectService) Detail(userId, projectId uint) (*domain.Project, error) {
	ok, err := s.storage.Project.IsOwned(userId, projectId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	project, err := s.storage.Project.FindById(projectId)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) Create(userId uint, input *domain.CreateProjectInput) (*domain.Project, error) {
	project := &domain.Project{
		UserId:      userId,
		Name:        input.Name,
		Description: input.Description,
		Deadline:    input.Deadline,
	}

	projectId, err := s.storage.Project.Create(project)
	if err != nil {
		return nil, err
	}

	project, err = s.storage.Project.FindById(projectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	err = s.storage.Heatmap.AddActivity(userId)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	return project, nil
}

func (s *projectService) Update(userId, projectId uint, input *domain.UpdateProjectInput) (*domain.Project, error) {
	ok, err := s.storage.Project.IsOwned(userId, projectId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	err = s.storage.Project.Update(&domain.Project{
		Id:          projectId,
		Name:        input.Name,
		Description: input.Description,
		Archived:    input.Archived,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	project, err := s.storage.Project.FindById(projectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	err = s.storage.Heatmap.AddActivity(userId)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	return project, nil
}

func (s *projectService) Delete(userId, projectId uint) error {
	ok, err := s.storage.Project.IsOwned(userId, projectId)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrUserNotOwnedRecord
	}

	if err := s.storage.Project.Delete(projectId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrRecordNotFound
		}
		return err
	}

	err = s.storage.Heatmap.RemoveActivity(userId)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	return nil
}

func (s *projectService) GetStats(userId, projectId uint) (*domain.ProjectStats, error) {
	ok, err := s.storage.Project.IsOwned(userId, projectId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	return s.storage.Project.GetStats(projectId)
}

func (s *projectService) GetMetrics(userId, projectId uint) (*domain.ProjectMetrics, error) {
	ok, err := s.storage.Project.IsOwned(userId, projectId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	preMetrics, err := s.storage.Project.GetMetrics(projectId)
	if err != nil {
		return nil, err
	}

	metrics := &domain.ProjectMetrics{
		TotalTasks:  preMetrics.TotalTasks,
		DoneTasks:   preMetrics.DoneTasks,
		RemTasks:    preMetrics.RemTasks,
		DaysElapsed: preMetrics.DaysElapsed,
		DaysLeft:    preMetrics.DaysLeft,
	}

	if metrics.DaysElapsed == 0 {
		metrics.VReal = float64(metrics.DoneTasks)
	} else {
		metrics.VReal = float64(metrics.DoneTasks) / float64(metrics.DaysElapsed)
	}

	if metrics.DaysLeft == 0 {
		metrics.VReq = float64(metrics.RemTasks)
	} else {
		metrics.VReq = float64(metrics.RemTasks) / float64(metrics.DaysLeft)
	}

	if metrics.TotalTasks == 0 {
		metrics.PerceptionDone = 0
	} else {
		metrics.PerceptionDone = int((float64(metrics.DoneTasks) / float64(metrics.TotalTasks)) * 100)
	}

	if metrics.VReal > 0 {
		offset := math.Ceil(float64(metrics.RemTasks) / metrics.VReal)

		now := time.Now().UTC()
		now = now.Add(time.Duration(offset) * 24 * time.Hour)
		metrics.ProjectedFinishDate = now
	}

	project, err := s.storage.Project.FindById(projectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	var status string

	now := time.Now().UTC()

	if project.Deadline.Before(now) {
		status = "red"
	} else {
		if metrics.VReal > metrics.VReq {
			status = "green"
		} else if metrics.VReal < metrics.VReq && metrics.DaysLeft > 0 {
			status = "yellow"
		} else if metrics.DaysLeft == 0 || (metrics.VReal == 0 && metrics.RemTasks > 0) {
			status = "red"
		}
	}

	metrics.Status = status

	return metrics, nil
}

func (s *projectService) GetProgress(userId, projectId uint) ([]*domain.ProjectProgress, error) {
	ok, err := s.storage.Project.IsOwned(userId, projectId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrUserNotOwnedRecord
	}

	return s.storage.Project.GetProgress(projectId)
}
