package domain

import "time"

type Project struct {
	Id          uint       `json:"id"`
	UserId      uint       `json:"user_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
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
