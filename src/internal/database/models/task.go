package models

import "time"

type Task struct {
	Id          uint `gorm:"primaryKey"`
	ColumnId    uint
	Name        string `gorm:"type:varchar(25);not null"`
	Description string `gorm:"type:text"`
	Status      bool
	Deadline    *time.Time `gorm:"type:timestamptz"`
	CreatedAt   time.Time  `gorm:"type:timestamptz"`
	UpdatedAt   time.Time  `gorm:"type:timestamptz"`
}

func (t *Task) TableName() string {
	return "tasks"
}
