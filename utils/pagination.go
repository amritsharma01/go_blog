// utils/pagination.go

package utils

import (
	"github.com/labstack/echo/v4"
)

func PaginatedResponse(c echo.Context, status int, message string, data interface{}, page, limit int, total int64) error {
	totalPages := int((total + int64(limit) - 1) / int64(limit)) // round up

	response := echo.Map{
		"data":       data,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
	}

	return JSONResponse(c, status, message, response)
}
