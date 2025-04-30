package handlers

import (
	requestmodels "crud_api/request_models"
	responsemodels "crud_api/response_models"
	"crud_api/services"
	"crud_api/utils"
	"net/http"

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

	req.Sanitize()
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Missing required fields")
	}

	user := requestmodels.FromUserCreateRequest(req)
	if err := h.service.Register(&user); err != nil {
		return HandleAppError(c, err, "Unexpected error happened diring registering user")
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

	if req.Email == "" || req.Password == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Missing email or password")
	}

	user, token, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		return HandleAppError(c, err, "Unexpected error happened during login")
	}

	resp := responsemodels.NewLoginResponse(responsemodels.ToUserResponse(*user), token)
	return utils.JSONResponse(c, http.StatusOK, "Login successful", resp)
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
		return HandleAppError(c, err, "Unexpected error occoured during retrieving")
	}

	var response []responsemodels.UserResponse
	for _, u := range users {
		response = append(response, responsemodels.ToUserResponse(u))
	}

	return utils.JSONResponse(c, http.StatusOK, "Successfully retrieved users", response)
}
