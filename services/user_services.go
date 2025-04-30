package services

import (
	"crud_api/errors"
	"crud_api/models"
	"crud_api/repositories"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *models.User) error
	Authenticate(email, password string) (*models.User, string, error)
	GetAllUsers() ([]models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(user *models.User) error {
	existingUser, err := s.repo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.Conflict("User with this email already exists")
	}

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if errHash != nil {
		return errors.Internal("Password hashing failed", errHash)
	}
	user.Password = string(hashedPassword)

	if errCreate := s.repo.Create(user); errCreate != nil {
		return errCreate
	}
	return nil
}

func (s *userService) Authenticate(email, password string) (*models.User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.Unauthorized("Invalid email or password")
	}

	// JWT generation
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, "", errors.Internal("Failed to generate token", err)
	}

	return user, signedToken, nil
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}
