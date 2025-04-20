package domain

import "time"

type Project struct {
	Id          uint       `json:"id"`
	UserId      uint       `json:"-"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
}

type CreateProjectInput struct {
	Name        string     `json:"name" validate:"required,gte=3,lte=50"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateProjectInput struct {
	Name        string     `json:"name" validate:"required,gte=3,lte=50"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
}
