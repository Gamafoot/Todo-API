package domain

import "time"

type Project struct {
	Id          uint
	UserId      uint
	Name        string
	Description string
	Deadline    *time.Time
}

type CreateProjectInput struct {
	Name        string     `json:"name" binding:"required,min=3,max=50"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateProjectInput struct {
	Name        string     `json:"name" binding:"required,min=3,max=50"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
}
