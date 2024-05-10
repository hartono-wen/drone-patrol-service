package validator

import (
	"github.com/go-playground/validator/v10"
)

type RequestValidatorInterface interface {
	Validate(i interface{}) error
}

type (
	RequestValidator struct {
		validator *validator.Validate
	}
)

func (cv *RequestValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{
		validator: validator.New(),
	}
}
