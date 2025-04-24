package handlers

import (
	"crud_api/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

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
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request body"})
	}

	// Check if all required fields are filled
	if category.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Category Name is required"})
	}

	var existing models.Category
	if err := h.DB.Where("cname = ?", category.Name).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"message": "Category already exists"})
	}

	// Save post
	if err := h.DB.Create(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to add category"})
	}

	return c.JSON(http.StatusCreated, category)
}
