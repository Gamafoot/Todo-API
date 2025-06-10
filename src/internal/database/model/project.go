package model

import "time"

type Project struct {
	Id          uint       `gorm:"primaryKey"`
	UserId      uint       `gorm:"not null"`
	Name        string     `gorm:"type:varchar(50);not null"`
	Description *string    `gorm:"type:text"`
	Archived    *bool      `gorm:"default:false"`
	Deadline    *time.Time `gorm:"type:timestamptz"`
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   time.Time  `gorm:"not null"`
	User        User       `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (p Project) TableName() string {
	return "projects"
}
