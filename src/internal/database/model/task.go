package model

import "time"

type Task struct {
	Id          uint       `gorm:"primaryKey"`
	ColumnId    uint       `gorm:"not null"`
	Name        string     `gorm:"type:varchar(50);not null"`
	Description *string    `gorm:"type:text"`
	Status      *bool      `gorm:"default:false"`
	CompletedAt *time.Time `gorm:"default:null"`
	Archived    *bool      `gorm:"default:false"`
	Position    int
	Deadline    *time.Time `gorm:"type:timestamptz"`
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   time.Time  `gorm:"not null"`
	Column      Column     `gorm:"foreignKey:ColumnId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *Task) TableName() string {
	return "tasks"
}
