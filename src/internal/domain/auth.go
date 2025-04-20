package domain

import "time"

type Session struct {
	Id           uint      `json:"-"`
	UserId       uint      `json:"-"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterInput struct {
	Username   string `json:"username" validate:"required,gte=3,lte=25"`
	Password   string `json:"password" validate:"required,gte=8,lte=64"`
	RePassword string `json:"re_password" validate:"required,gte=8,lte=64"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
