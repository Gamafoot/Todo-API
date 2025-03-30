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
	ColumnId    uint      `json:"column_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	Deadline    time.Time `json:"timestamp"`
}

type UpdateTaskInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      bool   `json:"status,omitempty"`
	Deadline    string `json:"timestamp,omitempty"`
}
