package app

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (cv *Validator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
