package pkg

import (
	"fmt"
	"log"
	"runtime"
)

type CustomError struct {
	Reason     string `json:"reason"`
	StatusCode int    `json:"statusCode"`
	Location   string `json:"location"` // More flexible
}

// Error implements the error interface
func (e *CustomError) Error() string {
	return fmt.Sprintf("[%d] %s (at %s)", e.StatusCode, e.Reason, e.Location)
}

// New creates and automatically handles the error
// Returns the error so you can return it from functions
func New(reason string, statusCode int) *CustomError {
	// Get caller location for better debugging
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)

	err := &CustomError{
		Reason:     reason,
		StatusCode: statusCode,
		Location:   location,
	}

	// Automatic handling
	log.Printf("ERROR: %s", err.Error())
	// Could add metrics, sentry, etc. here

	return err
}
