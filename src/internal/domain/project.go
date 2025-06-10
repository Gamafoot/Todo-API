package domain

import "time"

type Project struct {
	Id          uint       `json:"id"`
	UserId      uint       `json:"-"`
	Name        string     `json:"name"`
	Archived    *bool      `json:"archived"`
	Description *string    `json:"description"`
	Deadline    *time.Time `json:"deadline"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateProjectInput struct {
	Name        string     `json:"name" validate:"required,gte=3,lte=50"`
	Description *string    `json:"description"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateProjectInput struct {
	Name        string     `json:"name" validate:"omitempty,gte=3,lte=50"`
	Description *string    `json:"description"`
	Archived    *bool      `json:"archived"`
	Deadline    *time.Time `json:"deadline"`
}

type ProjectStats struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Overdue   int `json:"overdue"`
}

var ProjectOrder = map[string]string{
	"updated_at": "updated_at",
	"created_at": "created_at",
}

type SearchProjectOptions struct {
	Pattern  string
	Order    string
	Archived bool
}

type PreProjectMetrics struct {
	TotalTasks  int `json:"total_tasks"`
	DoneTasks   int `json:"done_tasks"`
	RemTasks    int `json:"rem_tasks"`
	DaysElapsed int `json:"days_elapsed"`
	DaysLeft    int `json:"days_left"`
}

type ProjectMetrics struct {
	TotalTasks          int       `json:"total_tasks"`
	DoneTasks           int       `json:"done_tasks"`
	RemTasks            int       `json:"rem_tasks"`
	DaysElapsed         int       `json:"days_elapsed"`
	DaysLeft            int       `json:"days_left"`
	VReal               float64   `json:"v_real"`
	VReq                float64   `json:"v_req"`
	PerceptionDone      int       `json:"perception_done"`
	ProjectedFinishDate time.Time `json:"projected_finish_date"`
	Status              string    `json:"status"`
}

type ProjectProgress struct {
	Day   time.Time `json:"day"`
	Count int       `json:"count"`
}
