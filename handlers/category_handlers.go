package handlers

import (
	"crud_api/errors"
	requestmodels "crud_api/request_models"
	responsemodels "crud_api/response_models"
	"crud_api/services"

	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	service services.CategoryService
}

// NewCategoryHandler returns a new instance of CategoryHandler
func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// AddCategory godoc
// @Summary Add a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body requestmodels.CategoryRequest true "Category data"
// @Success 201 {object} responsemodels.JSONResponseStruct{data=responsemodels.CategoryResponse}
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /v1/categories [post]
func (h *CategoryHandler) AddCategory(c echo.Context) error {
	var req requestmodels.CategoryRequest

	// Bind incoming JSON request
	if err := c.Bind(&req); err != nil {
		return errors.HandleError(c,
			errors.BadRequest(
				"Invalid request body",
				"Failed to bind request body",
				err,
			),
			"",
		)
	}

	// Check if category name is provided
	if req.Name == "" {
		return errors.HandleError(c,
			errors.BadRequest(
				"Category name is required",
				"Client sent empty category name",
				nil,
			),
			"",
		)
	}

	// Convert to domain model
	category := requestmodels.FromCatRequest(req)

	// Add category through service layer
	if err := h.service.AddCategory(&category); err != nil {
		return errors.HandleError(c, err, "Failed to create category")
	}

	// Return success response with category details
	return responsemodels.JSONResponse(c, http.StatusCreated, "Category created successfully", responsemodels.ToCatResponse(category))
}

// ListCategories godoc
// @Summary List categories
// @Description Get a paginated list of categories
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} responsemodels.PaginatedResponse{data=[]responsemodels.CategoryResponse}
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /v1/categories [get]
func (h *CategoryHandler) ListCategories(c echo.Context) error {
	// Get pagination details
	p := responsemodels.GetPagination(c)

	// Fetch categories from service
	categories, total, err := h.service.GetCategories(p.Limit, p.Offset)
	if err != nil {
		return errors.HandleError(c, err, "Failed to retrieve categories")
	}

	// Convert categories to response models
	var response []responsemodels.CategoryResponse
	for _, cat := range categories {
		response = append(response, responsemodels.ToCatResponse(cat))
	}

	// Return paginated response
	paginated := responsemodels.NewPaginatedResponse(response, p.Page, p.Limit, total)
	return responsemodels.SendPaginatedResponse(c, http.StatusOK, "Categories retrieved successfully", paginated)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	// Extract category ID from URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.HandleError(c,
			errors.BadRequest(
				"Invalid category ID",
				"Failed to parse category ID as integer",
				err,
			),
			"",
		)
	}

	// Fetch category by ID from service layer
	cat, err := h.service.GetByID(uint(id))
	if err != nil {
		return errors.HandleError(c, err, "Failed to delete category")
	}

	// Attempt to delete category through service layer
	if err := h.service.DeleteCategory(cat); err != nil {
		return errors.HandleError(c, err, "Failed to delete category")
	}

	// Return success response
	return responsemodels.JSONResponse(c, http.StatusOK, "Category deleted successfully", nil)
}
