package domain

import "time"

type Subtask struct {
	Id        uint
	TaskId    uint
	Name      string
	Status    bool
	Deadline  *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateSubtaskInput struct {
	TaskId   uint       `json:"task_id" binding:"required"`
	Name     string     `json:"name" binding:"required,min=3,max=50"`
	Status   bool       `json:"status" binding:"required"`
	Deadline *time.Time `json:"deadline"`
}

type UpdateSubtaskInput struct {
	Name        string     `json:"name" binding:"min=3,max=50"`
	Description string     `json:"description"`
	Status      bool       `json:"status"`
	Deadline    *time.Time `json:"deadline"`
}
