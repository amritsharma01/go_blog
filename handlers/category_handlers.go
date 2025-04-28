package handlers

import (
	"crud_api/models"
	responsemodels "crud_api/response_models"
	"crud_api/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func toCatResponse(c models.Category) responsemodels.CategoryResponse {
	return responsemodels.CategoryResponse{
		ID:   c.ID,
		Name: c.Name,
	}
}

type CategoryHandler struct {
	DB *gorm.DB
}

// Constructor
func NewCategory(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

func (h *CategoryHandler) AddCategory(c echo.Context) error {
	var category models.Category

	// Bind the request body to the Post struct
	if err := c.Bind(&category); err != nil {
		// Log the error message
		c.Logger().Errorf("Error binding request body: %v", err)
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")

	}

	// Check if all required fields are filled
	if category.Name == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Category Name is required")

	}

	var existing models.Category
	if err := h.DB.Where("cname = ?", category.Name).First(&existing).Error; err == nil {
		return utils.ErrorResponse(c, http.StatusConflict, "Category already exists")
	}

	// Save post
	if err := h.DB.Create(&category).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to add category")
	}
	return utils.JSONResponse(c, http.StatusCreated, "Category Created Succesfully", toCatResponse(category))

}

func (h *CategoryHandler) ListCategories(c echo.Context) error {
	var categories []models.Category

	pageNum := 1
	limitNum := 10

	if p, err := strconv.Atoi(c.QueryParam("page")); err == nil && p > 0 {
		pageNum = p
	}
	if l, err := strconv.Atoi(c.QueryParam("limit")); err == nil && l > 0 {
		limitNum = l
	}

	offset := (pageNum - 1) * limitNum

	if err := h.DB.Limit(limitNum).Offset(offset).Find(&categories).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve categories")
	}

	var total int64
	if err := h.DB.Model(&models.Category{}).Count(&total).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to count categories")
	}

	var responseCategories []responsemodels.CategoryResponse
	for _, cat := range categories {
		responseCategories = append(responseCategories, toCatResponse(cat))
	}

	return utils.PaginatedResponse(c, http.StatusOK, "Categories retrieved successfully", responseCategories, pageNum, limitNum, total)
}
