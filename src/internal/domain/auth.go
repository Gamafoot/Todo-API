package domain

import "time"

type Session struct {
	Id           uint
	UserId       uint
	RefreshToken string
	ExpiresAt    time.Time
}

type LoginInput struct {
	Username string `json:"username" binding:"required,min=3,max=25"`
	Password string `json:"password" binding:"required,min=3,max=64"`
}

type RegisterInput struct {
	Username   string `json:"username" binding:"required,min=3,max=25"`
	Password   string `json:"password" binding:"required,min=8,max=64"`
	RePassword string `json:"re_password" binding:"required,min=8,max=64"`
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
