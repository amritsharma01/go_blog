package handlers

import (
	models "crud_api/models"
	responsemodels "crud_api/response_models"
	"crud_api/utils"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ToUserResponse(u models.User) responsemodels.UserResponse {
	return responsemodels.UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

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
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Request")
	}

	// Check email uniqueness else error
	var existing models.User
	if err := h.DB.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		return utils.ErrorResponse(c, http.StatusConflict, "Email Already Exists")
	}

	// Hash password else error
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
	}
	user.Password = string(hashedPassword)

	// Save user else error
	if err := h.DB.Create(&user).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
	}

	return utils.JSONResponse(c, http.StatusCreated, "User Succesfully Created", ToUserResponse(user))
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	var users []models.User

	// Retrieve all users from the database
	if err := h.DB.Find(&users).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve users")

	}

	var response []responsemodels.UserResponse

	//send only the required contentsfor a single user
	for _, u := range users {
		response = append(response, ToUserResponse(u))

	}

	return utils.JSONResponse(c, http.StatusOK, "Succesfully Retrieved all users", response)

}

func (h *UserHandler) Login(c echo.Context) error {
	jwtSecret := os.Getenv("jwtSecret")
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&creds); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request"})
	}

	var user models.User
	if err := h.DB.Where("email = ?", creds.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "User not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid credentials"})
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Could not generate token"})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": signedToken})
}
