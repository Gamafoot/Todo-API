package models

import "time"

type User struct {
	Id        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"type:varchar(13);unique;not null"`
	Password  string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null"`
}

func (u User) TableName() string {
	return "users"
}
