// errors/errors.go
package errors

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AppErrors struct {
	Code        int    `json:"-"`       // HTTP status code
	Message     string `json:"message"` // User-friendly message
	InternalMsg string `json:"-"`       // Internal message for logging
	Err         error  `json:"-"`       // Original error
}

func (e *AppErrors) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.InternalMsg, e.Err)
	}
	return e.InternalMsg
}

func (e *AppErrors) Unwrap() error {
	return e.Err
}

// ErrorResponse
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"error"`
}

// Constructor functions for different error types
func BadRequest(userMsg string, InternalMsg string, err ...error) *AppErrors {
	var originalErr error
	if len(err) > 0 {
		originalErr = err[0]
	}
	return &AppErrors{
		Code:        http.StatusBadRequest,
		Message:     userMsg,
		InternalMsg: InternalMsg,
		Err:         originalErr,
	}
}

func NotFound(userMsg string, InternalMsg string, err ...error) *AppErrors {
	var originalErr error
	if len(err) > 0 {
		originalErr = err[0]
	}
	return &AppErrors{
		Code:        http.StatusNotFound,
		Message:     userMsg,
		InternalMsg: InternalMsg,
		Err:         originalErr,
	}
}

func Conflict(userMsg string, InternalMsg string, err ...error) *AppErrors {
	var originalErr error
	if len(err) > 0 {
		originalErr = err[0]
	}
	return &AppErrors{
		Code:        http.StatusConflict,
		Message:     userMsg,
		InternalMsg: InternalMsg,
		Err:         originalErr,
	}
}

func Internal(userMsg string, InternalMsg string, err ...error) *AppErrors {
	var originalErr error
	if len(err) > 0 {
		originalErr = err[0]
	}
	return &AppErrors{
		Code:        http.StatusInternalServerError,
		Message:     userMsg,
		InternalMsg: InternalMsg,
		Err:         originalErr,
	}
}

func Forbidden(userMsg string, InternalMsg string, err ...error) *AppErrors {
	var originalErr error
	if len(err) > 0 {
		originalErr = err[0]
	}
	return &AppErrors{
		Code:        http.StatusForbidden,
		Message:     userMsg,
		InternalMsg: InternalMsg,
		Err:         originalErr,
	}
}

func HandleError(c echo.Context, err error, defaultUserMsg string) error {
	statusCode := http.StatusInternalServerError
	userMsg := defaultUserMsg
	if userMsg == "" {
		userMsg = "An unexpected error occurred"
	}
	if appErr, ok := err.(*AppErrors); ok {
		logErrorWithSeverity(appErr)

		statusCode = appErr.Code
		if appErr.Message != "" {
			userMsg = appErr.Message
		}
	} else {
		log.Printf("SEVERE: Unhandled error: %v", err)
	}

	return c.JSON(statusCode, ErrorResponse{
		Status:  statusCode,
		Message: userMsg,
	})
}

func logErrorWithSeverity(appErr *AppErrors) {
	prefix := "INFO"
	if appErr.Code >= 400 && appErr.Code < 500 {
		prefix = "WARNING"
	} else if appErr.Code >= 500 {
		prefix = "SEVERE"
	}

	if appErr.Err != nil {
		log.Printf("%s: %s - Original error: %v", prefix, appErr.InternalMsg, appErr.Err)
	} else {
		log.Printf("%s: %s", prefix, appErr.InternalMsg)
	}
}
