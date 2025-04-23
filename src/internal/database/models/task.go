package models

import "time"

type Task struct {
	Id          uint       `gorm:"primaryKey"`
	ColumnId    uint       `gorm:"not null"`
	Name        string     `gorm:"type:varchar(25);not null"`
	Description string     `gorm:"type:text"`
	Status      bool       `gorm:"default:false"`
	Archived    bool       `gorm:"default:false"`
	Deadline    *time.Time `gorm:"type:timestamptz"`
	CreatedAt   time.Time  `gorm:"type:timestamptz;not null"`
	UpdatedAt   time.Time  `gorm:"type:timestamptz;not null"`
	Column      Column     `gorm:"foreignKey:ColumnId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *Task) TableName() string {
	return "tasks"
}
