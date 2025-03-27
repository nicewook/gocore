// Package valiecho provides a validator for Echo framework using go-playground/validator
package validatorutil

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator implements Echo's Validator interface using go-playground/validator
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new instance of CustomValidator
func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

// Validate implements the Echo Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
