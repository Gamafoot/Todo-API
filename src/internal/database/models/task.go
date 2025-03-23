package models

import "time"

type Task struct {
	Id          uint `gorm:"primaryKey"`
	ProjectId   uint
	Name        string `gorm:"type:varchar(25);not null"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(25)"`
	Deadline    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Project     Project `gorm:"foreignKey:ProjectId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *Task) TableName() string {
	return "tasks"
}
