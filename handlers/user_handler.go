package handlers

import (
	requestmodels "crud_api/request_models"
	responsemodels "crud_api/response_models"
	"crud_api/services"
	"crud_api/utils"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body requestmodels.CreateUserRequest true "User registration data"
// @Success 201 {object} utils.JSONResponseStruct{data=responsemodels.UserResponse}
// @Failure 400 {object} utils.ErrorResponseStruct
// @Failure 409 {object} utils.ErrorResponseStruct
// @Failure 500 {object} utils.ErrorResponseStruct
// @Router /auth/register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var req requestmodels.CreateUserRequest

	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	user := requestmodels.FromUserCreateRequest(req)

	if err := h.service.Register(&user); err != nil {
		if err == services.ErrEmailAlreadyExists {
			return utils.ErrorResponse(c, http.StatusConflict, "Email already exists")
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user")
	}

	return utils.JSONResponse(c, http.StatusCreated, "User created successfully", responsemodels.ToUserResponse(user))
}

// Login godoc
// @Summary Authenticate user
// @Description Login with email and password to get JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body requestmodels.LoginRequest true "Login credentials"
// @Success 200 {object} utils.JSONResponseStruct{data=responsemodels.LoginResponse}
// @Failure 400 {object} utils.ErrorResponseStruct
// @Failure 401 {object} utils.ErrorResponseStruct
// @Failure 500 {object} utils.ErrorResponseStruct
// @Router /auth/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var req requestmodels.LoginRequest

	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	loginData := requestmodels.FromUserLoginRequest(req)

	dbUser, err := h.service.Login(loginData.Email, loginData.Password)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"user_id": dbUser.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
	}

	loginResp := responsemodels.NewLoginResponse(responsemodels.ToUserResponse(*dbUser), signedToken)
	return utils.JSONResponse(c, http.StatusOK, "Login successful", loginResp)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve list of all users (requires admin JWT)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.PaginatedResponse{data=[]responsemodels.UserResponse}
// @Failure 401 {object} utils.ErrorResponseStruct
// @Failure 403 {object} utils.ErrorResponseStruct
// @Failure 500 {object} utils.ErrorResponseStruct
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve users")
	}

	var response []responsemodels.UserResponse
	for _, u := range users {
		response = append(response, responsemodels.ToUserResponse(u))
	}

	return utils.JSONResponse(c, http.StatusOK, "Successfully retrieved users", response)
}
