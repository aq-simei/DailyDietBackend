package errors

import "fmt"

type ErrorType string

const (
	NotFound     ErrorType = "NOT_FOUND"
	Internal     ErrorType = "INTERNAL"
	Invalid      ErrorType = "INVALID"
	Unauthorized ErrorType = "UNAUTHORIZED"
	Forbidden    ErrorType = "FORBIDDEN"
)

type CustomError struct {
	Type    ErrorType
	Message string
	Err     error
}

type DatabaseError struct {
	Type ErrorType
	Err  error
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func NewError(errType ErrorType, message string, err error) *CustomError {
	return &CustomError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}
