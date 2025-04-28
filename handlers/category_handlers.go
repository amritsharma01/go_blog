package handlers

import (
	"crud_api/models"
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

// Constructor
func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) AddCategory(c echo.Context) error {
	var req requestmodels.CategoryRequest

	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// Validate manually if needed (or use echo middleware later)
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

	category := models.Category{
		Name: req.Name,
	}

	if err := h.service.AddCategory(&category); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category")
	}

	return utils.JSONResponse(c, http.StatusCreated, "Category created successfully", responsemodels.ToCatResponse(category))
}

func (h *CategoryHandler) ListCategories(c echo.Context) error {
	pageNum := 1
	limitNum := 10

	if p, err := strconv.Atoi(c.QueryParam("page")); err == nil && p > 0 {
		pageNum = p
	}
	if l, err := strconv.Atoi(c.QueryParam("limit")); err == nil && l > 0 {
		limitNum = l
	}
	offset := (pageNum - 1) * limitNum

	categories, total, err := h.service.GetCategories(limitNum, offset)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve categories")
	}

	var responseCategories []responsemodels.CategoryResponse
	for _, cat := range categories {
		responseCategories = append(responseCategories, responsemodels.ToCatResponse(cat))
	}

	return utils.PaginatedResponse(c, http.StatusOK, "Categories retrieved successfully", responseCategories, pageNum, limitNum, total)
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	category_id, _ := strconv.Atoi(c.Param("id"))

	cat, err := h.service.GetByID(uint(category_id))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Category Not Found")
	}

	if err := h.service.DeleteCategory(cat); err != nil {
		return utils.ErrorResponse(c, http.StatusForbidden, "Unable to delete category")
	}

	return utils.JSONResponse(c, http.StatusOK, "Category Deleted Sccesfully", nil)
}
