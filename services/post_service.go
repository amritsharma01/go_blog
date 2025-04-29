package services

import (
	"crud_api/models"
	"crud_api/repositories"
	"errors"
)

type PostService interface {
	Create(post *models.Post) error
	GetByID(id uint) (*models.Post, error)
	GetAll(search, categoryID, authorID string, offset, limit int) ([]models.Post, int64, error)
	GetByAuthorID(authorID uint, offset, limit int) ([]models.Post, int64, error)
	Update(post *models.Post) error
	Delete(post *models.Post, userID uint) error
}

type postService struct {
	repo repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) PostService {
	return &postService{repo}
}

func (s *postService) Create(post *models.Post) error {
	if existing, err := s.repo.FindDuplicate(post.Title, post.AuthorID); err == nil && existing != nil {
		return ErrPostAlreadyExists
	}
	return s.repo.Create(post)
}

func (s *postService) GetByID(id uint) (*models.Post, error) {
	return s.repo.FindByID(id)
}

func (s *postService) GetAll(search, categoryID, authorID string, offset, limit int) ([]models.Post, int64, error) {
	return s.repo.FindAll(search, categoryID, authorID, offset, limit)
}

func (s *postService) GetByAuthorID(authorID uint, offset, limit int) ([]models.Post, int64, error) {
	return s.repo.FindByAuthorID(authorID, offset, limit)
}

func (s *postService) Update(post *models.Post) error {
	return s.repo.Update(post)
}

func (s *postService) Delete(post *models.Post, userID uint) error {
	if post.AuthorID != userID {
		return ErrUnauthorizedAction
	}
	return s.repo.Delete(post)
}

// Custom errors
var (
	ErrPostAlreadyExists  = errors.New("post already exists")
	ErrUnauthorizedAction = errors.New("you are not authorized to perform this action")
)

// struct AppError{
// 	message string
// 	status int
// }

// func ServerError(message string) error {
// 	retu
// }
