package handlers

import (
	requestmodels "crud_api/request_models"
	responsemodels "crud_api/response_models"
	"crud_api/services"
	"crud_api/utils"
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
// @Success 201 {object} utils.JSONResponseStruct{data=responsemodels.CategoryResponse}
// @Failure 400 {object} utils.ErrorResponseStruct
// @Failure 409 {object} utils.ErrorResponseStruct
// @Failure 500 {object} utils.ErrorResponseStruct
// @Router /categories [post]
func (h *CategoryHandler) AddCategory(c echo.Context) error {
	var req requestmodels.CategoryRequest

	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Category name is required")
	}

	exists, err := h.service.CategoryExists(req.Name)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Error checking category existence")
	}
	if exists {
		return utils.ErrorResponse(c, http.StatusConflict, "Category already exists")
	}

	category := requestmodels.FromCatRequest(req)

	if err := h.service.AddCategory(&category); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category")
	}

	return utils.JSONResponse(c, http.StatusCreated, "Category created successfully", responsemodels.ToCatResponse(category))
}

// ListCategories godoc
// @Summary List categories
// @Description Get a paginated list of categories
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]responsemodels.CategoryResponse}
// @Failure 500 {object} utils.ErrorResponseStruct
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c echo.Context) error {
	p := utils.GetPagination(c)

	categories, total, err := h.service.GetCategories(p.Limit, p.Offset)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve categories")
	}

	var response []responsemodels.CategoryResponse
	for _, cat := range categories {
		response = append(response, responsemodels.ToCatResponse(cat))
	}

	paginated := utils.NewPaginatedResponse(response, p.Page, p.Limit, total)
	return utils.SendPaginatedResponse(c, http.StatusOK, "Categories retrieved successfully", paginated)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} utils.JSONResponseStruct
// @Failure 403 {object} utils.ErrorResponseStruct
// @Failure 404 {object} utils.ErrorResponseStruct
// @Failure 500 {object} utils.ErrorResponseStruct
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	cat, err := h.service.GetByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
	}

	if err := h.service.DeleteCategory(cat); err != nil {
		return utils.ErrorResponse(c, http.StatusForbidden, "Unable to delete category")
	}

	return utils.JSONResponse(c, http.StatusOK, "Category deleted successfully", nil)
}
