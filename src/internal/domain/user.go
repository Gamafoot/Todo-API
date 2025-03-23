package domain

import "time"

type User struct {
	Id        uint
	Username  string
	Password  string
	CreatedAt time.Time
}
