package domain

import (
	"errors"
)

var (
	ErrUsernameIsOccupied     = errors.New("username is occupied")
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrPasswordsDontMatch     = errors.New("passwords don't match")
	ErrReshreshTokenNotFound  = errors.New("refresh token not found")
	ErrInvalidDeadlineFormat  = errors.New("invalid deadline format")
	ErrTaskNotFound           = errors.New("task not found")
)
