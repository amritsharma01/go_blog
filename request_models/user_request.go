package requestmodels

import (
	"crud_api/models"
	"strings"
)

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=1"`
	Email    string `json:"email" validate:"required,email,min=1"`
	Password string `json:"password" validate:"required,min=6"`
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

func (r *CreateUserRequest) Sanitize() {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)

}

func (r *LoginRequest) Sanitize() {

	r.Email = strings.TrimSpace(r.Email)

}
