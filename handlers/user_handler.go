package handlers

import (
	"crud_api/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

// Constructor
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user models.User

	// Try to bind the request body to the User struct else error
	if err := c.Bind(&user); err != nil {
		// Log the error message
		c.Logger().Errorf("Error binding request body: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request body"})
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "All fields are required"})
	}

	// Check email uniqueness else error
	var existing models.User
	if err := h.DB.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"message": "Email already exists"})
	}

	// Hash password else error
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to hash password"})
	}
	user.Password = string(hashedPassword)

	// Save user else error
	if err := h.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to create user"})
	}

	user.Password = "" // reurning empty string as pass
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	var users []models.User

	// Retrieve all users from the database
	if err := h.DB.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve users"})
	}

	// Hide passwords before sending the response
	for i := range users {
		users[i].Password = ""
	}

	return c.JSON(http.StatusOK, users)
}
