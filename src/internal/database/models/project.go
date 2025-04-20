package models

import "time"

type Project struct {
	Id          uint `gorm:"primaryKey"`
	UserId      uint
	Name        string     `gorm:"type:varchar(50);not null"`
	Description string     `gorm:"type:text"`
	Deadline    *time.Time `gorm:"type:timestamptz"`
}

func (p Project) TableName() string {
	return "projects"
}
