package domain

import "time"

type Session struct {
	Id           uint
	UserId       uint
	RefreshToken string
	ExpiresAt    time.Time
}
