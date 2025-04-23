package models

import "time"

type Session struct {
	Id           uint      `gorm:"primaryKey"`
	UserId       uint      `gorm:"not null"`
	RefreshToken string    `gorm:"type:varchar(255);not null"`
	ExpiresAt    time.Time `gorm:"not null"`
	User         User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
