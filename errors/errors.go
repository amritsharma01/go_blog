package errors

import "fmt"

type ErrorCode string

// AppError is a custom error for structured error handling
type AppError struct {
	Code    int    // HTTP status code
	Message string // User-facing error message
	Err     error  // Root/internal error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Code: %d, Message: %s, Details: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// Constructors

func New(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func NewWithMessage(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Common HTTP error constructors

func BadRequest(message string) *AppError {
	return NewWithMessage(400, message)
}

func Unauthorized(message string) *AppError {
	return NewWithMessage(401, message)
}

func Forbidden(message string) *AppError {
	return NewWithMessage(403, message)
}

func NotFound(message string) *AppError {
	return NewWithMessage(404, message)
}

func Conflict(message string) *AppError {
	return NewWithMessage(409, message)
}

func Internal(message string, err error) *AppError {
	return New(500, message, err)
}
