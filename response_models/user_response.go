package responsemodels

import "crud_api/models"

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ToUserResponse(u models.User) UserResponse {
	return UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

func NewLoginResponse(user UserResponse, token string) LoginResponse {
	return LoginResponse{
		User:  user,
		Token: token,
	}
}
