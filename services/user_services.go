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
	GetByID(id uint) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(user *models.User) error {
	existingUser, err := s.repo.FindByEmail(user.Email)
	if err != nil {
		// Allow only "not found" errors to proceed with creation
		if appErr, ok := err.(*errors.AppErrors); ok && appErr.Code == 404 {
			hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if errHash != nil {
				return errors.Internal("Failed to register the user", "Error hashing password", errHash)
			}
			user.Password = string(hashedPassword)
			if createdErr := s.repo.Create(user); createdErr != nil {
				return createdErr
			}
			return nil
		}
		return err
	}
	if existingUser != nil {
		return errors.Conflict("Failed to login with the given crrdentials", "Attempted to create a duplicate user")
	}

	return nil
}

func (s *userService) Authenticate(email, password string) (*models.User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.Internal("Invalid email or password", "User tried logging in with invalid email and password", err)
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
		return nil, "", err
	}

	return user, signedToken, nil
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}
