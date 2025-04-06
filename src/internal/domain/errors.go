package domain

import (
	pkgErrors "github.com/pkg/errors"
)

var (
	ErrNotDigit               = pkgErrors.New("not digit")
	ErrNotPositiveDigit       = pkgErrors.New("digit must be positive")
	ErrUsernameIsOccupied     = pkgErrors.New("username is occupied")
	ErrInvalidLoginOrPassword = pkgErrors.New("invalid login or password")
	ErrPasswordsDontMatch     = pkgErrors.New("passwords don't match")
	ErrReshreshTokenNotFound  = pkgErrors.New("refresh token not found")
	ErrRecordNotFound         = pkgErrors.New("record not found")
	ErrUserNotOwnedRecord     = pkgErrors.New("the user is not owner of the record")
)
