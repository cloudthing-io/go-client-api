package api

import (
	"fmt"
)

type ApiError struct {
	StatusCode int
	Message string
}

func (a ApiError) Error() string {
	return fmt.Sprintf("Status code: %d, error message: %s", a.StatusCode, a.Message)
}
