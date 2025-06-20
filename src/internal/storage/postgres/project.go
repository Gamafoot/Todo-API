package postgres

import (
	"root/internal/database/model"
	"root/internal/domain"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type projectStorage struct {
	db *gorm.DB
}

func NewProjectStorage(db *gorm.DB) *projectStorage {
	return &projectStorage{db: db}
}

func (s *projectStorage) FindAll(userId uint, options *domain.SearchProjectOptions, page, limit int) ([]*domain.Project, int, error) {
	offset := (page - 1) * limit

	// Субзапрос для подсчета задач и определения статуса
	taskStatsSubquery := s.db.Model(&model.Task{}).
		Select(`
            columns.project_id,
            COUNT(tasks.id) as task_count,
            BOOL_AND(COALESCE(tasks.status, false)) as all_completed
        `).
		Joins("JOIN columns ON columns.id = tasks.column_id").
		Where("tasks.archived IS NULL OR tasks.archived = false").
		Group("columns.project_id")

	// Субзапрос для определения финального статуса
	statusSubquery := s.db.Table("(?) as task_stats", taskStatsSubquery).
		Select(`
            project_id,
            CASE
                WHEN task_count = 0 THEN false
                ELSE all_completed
            END as computed_status
        `)

	// Базовый запрос
	baseQuery := s.db.Model(&model.Project{}).
		Select("projects.*, COALESCE(status_subquery.computed_status, false) as status").
		Joins("LEFT JOIN (?) as status_subquery ON projects.id = status_subquery.project_id", statusSubquery).
		Where("projects.user_id = ?", userId).
		Where("projects.archived = ?", options.Archived)

	if len(options.Pattern) > 0 {
		pattern := "%" + options.Pattern + "%"
		baseQuery = baseQuery.Where(
			"projects.name LIKE ? OR projects.description LIKE ?",
			pattern,
			pattern,
		)
	}

	// Подсчет общего количества
	var count int64
	countQuery := baseQuery.Session(&gorm.Session{})
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, pkgErrors.WithStack(err)
	}

	// Получение данных с пагинацией
	var projects []*model.Project
	findQuery := baseQuery.Session(&gorm.Session{})
	findQuery = findQuery.Offset(offset).Limit(limit).Order(options.Order + " DESC")
	if err := findQuery.Find(&projects).Error; err != nil {
		return nil, 0, pkgErrors.WithStack(err)
	}

	resultProjects := make([]*domain.Project, len(projects))
	for i, project := range projects {
		resultProjects[i] = toDomainProject(project)
	}

	return resultProjects, int(count), nil
}

func (s *projectStorage) FindById(id uint) (*domain.Project, error) {
	// Субзапрос для подсчета задач и определения статуса
	taskStatsSubquery := s.db.Model(&model.Task{}).
		Select(`
            columns.project_id,
            COUNT(tasks.id) as task_count,
            BOOL_AND(COALESCE(tasks.status, false)) as all_completed
        `).
		Joins("JOIN columns ON columns.id = tasks.column_id").
		Where("tasks.archived IS NULL OR tasks.archived = false").
		Group("columns.project_id")

	// Субзапрос для определения финального статуса
	statusSubquery := s.db.Table("(?) as task_stats", taskStatsSubquery).
		Select(`
            project_id,
            CASE 
                WHEN task_count = 0 THEN false
                ELSE all_completed
            END as computed_status
        `)

	// Основной запрос с добавлением вычисляемого поля status
	project := new(model.Project)
	if err := s.db.Model(&model.Project{}).
		Select("projects.*, COALESCE(status_subquery.computed_status, false) as status").
		Joins("LEFT JOIN (?) as status_subquery ON projects.id = status_subquery.project_id", statusSubquery).
		Where("projects.id = ?", id).
		First(project).Error; err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return toDomainProject(project), nil
}

func (s *projectStorage) Create(project *domain.Project) (uint, error) {
	modelProject := toModelProject(project)
	if err := s.db.Create(modelProject).Error; err != nil {
		return 0, pkgErrors.WithStack(err)
	}
	return modelProject.Id, nil
}

func (s *projectStorage) Update(project *domain.Project) error {
	if err := s.db.Model(model.Project{}).Where("id = ?", project.Id).Updates(project).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *projectStorage) Delete(id uint) error {
	if err := s.db.Delete(&model.Project{Id: id}).Error; err != nil {
		return pkgErrors.WithStack(err)
	}
	return nil
}

func (s *projectStorage) IsOwned(userId, projectId uint) (bool, error) {
	var isOwned bool

	err := s.db.Raw("SELECT is_owned_project(?, ?)", userId, projectId).Scan(&isOwned).Error
	if err != nil {
		return false, pkgErrors.WithStack(err)
	}

	return isOwned, nil
}

func (s *projectStorage) GetStats(projectId uint) (*domain.ProjectStats, error) {
	var (
		total     int64
		completed int64
		overdue   int64
	)

	query := s.db.Model(model.Task{}).Joins("INNER JOIN columns ON columns.id = tasks.column_id")
	err := query.Where("columns.project_id = ?", projectId).Count(&total).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	query = s.db.Model(model.Task{}).Joins("INNER JOIN columns ON columns.id = tasks.column_id")
	err = query.Where("columns.project_id = ? AND tasks.status = true", projectId).Count(&completed).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	query = s.db.Model(model.Task{}).Joins("INNER JOIN columns ON columns.id = tasks.column_id")
	err = query.Where("columns.project_id = ? AND tasks.deadline < (now() AT TIME ZONE 'UTC')", projectId).Count(&overdue).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return &domain.ProjectStats{
		Total:     int(total),
		Completed: int(completed),
		Overdue:   int(overdue),
	}, nil
}

func (s *projectStorage) GetMetrics(projectId uint) (*domain.PreProjectMetrics, error) {
	metrics := new(domain.PreProjectMetrics)

	err := s.db.Raw("SELECT * FROM project_metrics(?)", projectId).Scan(&metrics).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return metrics, nil
}

func (s *projectStorage) GetProgress(projectId uint) ([]*domain.ProjectProgress, error) {
	progress := make([]*domain.ProjectProgress, 0)

	err := s.db.Raw("SELECT * FROM project_progress(?)", projectId).Scan(&progress).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return progress, nil
}

func toDomainProject(project *model.Project) *domain.Project {
	return &domain.Project{
		Id:          project.Id,
		UserId:      project.UserId,
		Name:        project.Name,
		Description: project.Description,
		Deadline:    project.Deadline,
		Archived:    project.Archived,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
		Status:      project.Status,
	}
}

func toModelProject(project *domain.Project) *model.Project {
	return &model.Project{
		Id:          project.Id,
		UserId:      project.UserId,
		Name:        project.Name,
		Description: project.Description,
		Deadline:    project.Deadline,
		Archived:    project.Archived,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
		Status:      project.Status,
	}
}

func (s *projectStorage) setConditionForFindAll(db *gorm.DB) {

}
