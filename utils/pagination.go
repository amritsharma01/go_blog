// utils/pagination.go

package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"totalPages"`
}

// NewPaginatedResponse constructs a PaginatedResponse and returns it
func NewPaginatedResponse(data interface{}, page, limit int, total int64) PaginatedResponse {
	totalPages := int((total + int64(limit) - 1) / int64(limit)) // round up

	return PaginatedResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// SendPaginatedResponse is used to send a JSON response with pagination
func SendPaginatedResponse(c echo.Context, status int, message string, response PaginatedResponse) error {
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
