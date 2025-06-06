package model

import "time"

type Subtask struct {
	Id        uint   `gorm:"primaryKey"`
	TaskId    uint   `gorm:"not null"`
	Name      string `gorm:"type:varchar(25);not null"`
	Status    *bool  `gorm:"default:false"`
	Archived  *bool  `gorm:"default:false"`
	Position  int
	Deadline  *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	Task      Task       `gorm:"foreignKey:TaskId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *Subtask) TableName() string {
	return "subtasks"
}
