package configs

import "errors"

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

func NewValidationError(msg string) error {
	return &ValidationError{Message: msg}
}
