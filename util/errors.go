package util

import "fmt"

type Code struct {
	StatusCode  string
	Description string
}

type APIError struct {
	Code
}

func NewAPIError(statusCode string, description string) error {
	return &APIError{
		Code: Code{
			StatusCode:  statusCode,
			Description: description,
		},
	}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.StatusCode, e.Description)
}
