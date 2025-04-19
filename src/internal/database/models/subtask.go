package models

import "time"

type Subtask struct {
	Id        uint `gorm:"primaryKey"`
	TaskId    uint
	Name      string `gorm:"type:varchar(25);not null"`
	Status    bool
	Deadline  *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"type:timestamptz"`
	UpdatedAt time.Time  `gorm:"type:timestamptz"`
}

func (t *Subtask) TableName() string {
	return "subtasks"
}
