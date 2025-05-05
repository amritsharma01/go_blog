package responsemodels

import "github.com/labstack/echo/v4"

type JSONResponseStruct struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponseStruct struct {
	Error string `json:"error"`
}

func JSONResponse(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, JSONResponseStruct{
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, ErrorResponseStruct{
		Error: message,
	})
}
