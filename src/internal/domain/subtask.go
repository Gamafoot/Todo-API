package domain

import "time"

type Subtask struct {
	Id        uint      `json:"id"`
	TaskId    uint      `json:"task_id"`
	Name      string    `json:"name"`
	Status    bool      `json:"status"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSubtaskInput struct {
	TaskId uint   `json:"task_id" validate:"required"`
	Name   string `json:"name" validate:"required,gte=3,lte=50"`
}

type UpdateSubtaskInput struct {
	Name        string `json:"name" validate:"omitempty,gte=3,lte=50"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	Archived    bool   `json:"archived"`
}
