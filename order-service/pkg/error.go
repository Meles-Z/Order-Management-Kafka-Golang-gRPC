package pkg

import "fmt"

type CustomError struct {
	Reason     string `json:"reason"`
	StatusCode int    `json:"statusCode"`
	HappenIn   string `json:"happenIn"`
}

// Error implements the error interface
func (e CustomError) Error() string {
	return fmt.Sprintf("error: %s | status: %d | in: %s", e.Reason, e.StatusCode, e.HappenIn)
}

// Constructor
func NewCustomError(reason string, statusCode int, happenIn string) error {
	return CustomError{
		Reason:     reason,
		StatusCode: statusCode,
		HappenIn:   happenIn,
	}
}
