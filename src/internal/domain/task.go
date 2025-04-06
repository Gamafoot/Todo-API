package domain

import "time"

type Task struct {
	Id          uint
	ColumnId    uint
	Name        string
	Description string
	Status      bool
	Deadline    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateTaskInput struct {
	ColumnId    uint       `json:"column_id" binding:"required"`
	Name        string     `json:"name" binding:"required,min=3,max=50"`
	Description string     `json:"description" binding:"required"`
	Status      bool       `json:"status" binding:"required"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateTaskInput struct {
	Name        string     `json:"name" binding:"min=3,max=50"`
	Description string     `json:"description"`
	Status      bool       `json:"status"`
	Deadline    *time.Time `json:"deadline"`
}
