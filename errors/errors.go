package errors

import "fmt"

// Custom error type to represent a business logic error
type AppError struct {
	Code    int    // HTTP status code or business error code
	Message string // Error message
	Err     error  // The underlying error (optional)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Code: %d, Message: %s, Details: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// Helper function to create a new error
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Helper function to create a new error with no underlying error
func NewWithMessage(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
