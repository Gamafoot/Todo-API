package domain

import "time"

type Task struct {
	Id          uint       `json:"id"`
	ColumnId    uint       `json:"column_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      *bool      `json:"status"`
	Archived    *bool      `json:"archived"`
	Deadline    *time.Time `json:"deadline"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTaskInput struct {
	ColumnId    uint       `json:"column_id" validate:"required"`
	Name        string     `json:"name" validate:"required,gte=3,lte=50"`
	Description string     `json:"description" validate:"required"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateTaskInput struct {
	ColumnId    uint       `json:"column_id"`
	Name        string     `json:"name" validate:"omitempty,gte=3,lte=50"`
	Description string     `json:"description"`
	Status      *bool      `json:"status"`
	Archived    *bool      `json:"archived"`
	Deadline    *time.Time `json:"deadline"`
}
