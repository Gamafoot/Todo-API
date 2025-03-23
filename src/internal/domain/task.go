package domain

import "time"

type Task struct {
	Id          uint
	ProjectId   uint
	Name        string
	Description string
	Status      string
	Deadline    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
