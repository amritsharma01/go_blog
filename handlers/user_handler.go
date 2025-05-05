package handlers

import (
	"crud_api/errors"
	requestmodels "crud_api/request_models"
	responsemodels "crud_api/response_models"
	"crud_api/services"

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
// @Success 201 {object} responsemodels.JSONResponseStruct{data=responsemodels.UserResponse}
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /v1/auth/register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var req requestmodels.CreateUserRequest

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
	req.Sanitize()
	// Check if category name is provided
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return errors.HandleError(c,
			errors.BadRequest(
				"Username, password and email is required",
				"Client sent emmpy name, email or password",
				nil,
			),
			"",
		)
	}

	user := requestmodels.FromUserCreateRequest(req)
	if err := h.service.Register(&user); err != nil {
		return errors.HandleError(c, err, "")
	}

	return responsemodels.JSONResponse(c, http.StatusCreated, "User created successfully", responsemodels.ToUserResponse(user))
}

// Login godoc
// @Summary Authenticate user
// @Description Login with email and password to get JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body requestmodels.LoginRequest true "Login credentials"
// @Success 200 {object} responsemodels.JSONResponseStruct{data=responsemodels.LoginResponse{user=responsemodels.UserResponse}}
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /v1/auth/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var req requestmodels.LoginRequest

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

	req.Sanitize()

	// Check if category name is provided
	if req.Email == "" || req.Password == "" {
		return errors.HandleError(c,
			errors.BadRequest(
				"Email or password is required",
				"Client sent empty email or password",
				nil,
			),
			"",
		)
	}

	user, token, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		return errors.HandleError(c, err, "")
	}

	resp := responsemodels.NewLoginResponse(responsemodels.ToUserResponse(*user), token)
	return responsemodels.JSONResponse(c, http.StatusOK, "Login successful", resp)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve list of all users (requires admin JWT)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} responsemodels.PaginatedResponse{data=[]responsemodels.UserResponse}
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /v1/users [get]
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return errors.HandleError(c, err, "")
	}

	var response []responsemodels.UserResponse
	for _, u := range users {
		response = append(response, responsemodels.ToUserResponse(u))
	}

	return responsemodels.JSONResponse(c, http.StatusOK, "Successfully retrieved users", response)
}
