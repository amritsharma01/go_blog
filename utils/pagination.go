// utils/pagination.go

package utils

import (
	"strconv"

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

type Pagination struct {
	Page   int
	Limit  int
	Offset int
}

// GetPagination extracts pagination details from query params
func GetPagination(c echo.Context) Pagination {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	return Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}
