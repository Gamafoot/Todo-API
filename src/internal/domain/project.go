package domain

import "time"

type Project struct {
	Id          uint       `json:"id"`
	UserId      uint       `json:"-"`
	Name        string     `json:"name"`
	Archived    *bool      `json:"archived"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateProjectInput struct {
	Name        string     `json:"name" validate:"required,gte=3,lte=50"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateProjectInput struct {
	Name        string     `json:"name" validate:"required,gte=3,lte=50"`
	Description string     `json:"description"`
	Archived    *bool      `json:"archived"`
	Deadline    *time.Time `json:"deadline"`
}
