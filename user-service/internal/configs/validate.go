package configs

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate performs struct validation and returns well-formatted errors
func Validate(structure any) error {
	if err := validator.New().Struct(structure); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errs []string
			for _, e := range validationErrors {
				msg := strings.ToLower(e.Field()) + " must be " + e.Tag()
				if e.Param() != "" {
					msg += " " + e.Param()
				}
				errs = append(errs, msg)
			}
			return &ValidationError{Message: strings.Join(errs, "; ")}
		}
		return err // Return original error if not ValidationErrors
	}
	return nil
}

// ValidationError is a clean wrapper for validation failures
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return "validation error: " + e.Message
}
