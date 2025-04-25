package utils

import (
	"github.com/labstack/echo/v4"
)

func JSONResponse(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, echo.Map{
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, echo.Map{
		"error": message,
	})
}
