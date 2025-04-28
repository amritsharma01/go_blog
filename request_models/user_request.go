package requestmodels

import "crud_api/models"

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func FromUserCreateRequest(u CreateUserRequest) models.User {
	return models.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

func FromUserLoginRequest(u LoginRequest) models.User {
	return models.User{

		Email:    u.Email,
		Password: u.Password,
	}
}
