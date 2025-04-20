package domain

import "time"

type Subtask struct {
	Id        uint      `json:"id"`
	TaskId    uint      `json:"task_id"`
	Name      string    `json:"name"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
