package services

import (
	"crud_api/errors"
	"crud_api/models"
	"crud_api/repositories"
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
	existing, err := s.repo.FindDuplicate(post.Title, post.AuthorID)
	if err != nil {
		return err // propagate repository error with full context
	}
	if existing != nil {
		return errors.NewWithMessage(409, "Post already exists")
	}

	return s.repo.Create(post)
}

func (s *postService) GetByID(id uint) (*models.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.NewWithMessage(404, "Post not found")
	}
	return post, nil
}

func (s *postService) GetAll(search, categoryID, authorID string, offset, limit int) ([]models.Post, int64, error) {
	posts, count, err := s.repo.FindAll(search, categoryID, authorID, offset, limit)
	if err != nil {
		return nil, 0, err // keep raw repo error
	}
	return posts, count, nil
}

func (s *postService) GetByAuthorID(authorID uint, offset, limit int) ([]models.Post, int64, error) {
	posts, count, err := s.repo.FindByAuthorID(authorID, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return posts, count, nil
}

func (s *postService) Update(post *models.Post) error {
	return s.repo.Update(post) // propagate raw
}

func (s *postService) Delete(post *models.Post, userID uint) error {
	if post.AuthorID != userID {
		return errors.NewWithMessage(403, "You are not authorized to delete this post")
	}
	return s.repo.Delete(post)
}
