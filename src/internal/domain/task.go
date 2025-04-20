package domain

import "time"

type Task struct {
	Id          uint       `json:"id"`
	ColumnId    uint       `json:"column_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      bool       `json:"status"`
	Deadline    *time.Time `json:"deadline"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTaskInput struct {
	ColumnId    uint       `json:"column_id" binding:"required"`
	Name        string     `json:"name" binding:"required,min=3,max=50"`
	Description string     `json:"description" binding:"required"`
	Status      bool       `json:"status" binding:"required"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateTaskInput struct {
	ColumnId    uint       `json:"column_id"`
	Name        string     `json:"name" binding:"min=3,max=50"`
	Description string     `json:"description"`
	Status      bool       `json:"status"`
	Deadline    *time.Time `json:"deadline"`
}
