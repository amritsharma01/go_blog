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
	GetByAuthorID(authorID string, offset, limit int) ([]models.Post, int64, error)
	Update(post *models.Post, userID uint) error
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
		// Allow only "not found" errors to proceed with creation
		if appErr, ok := err.(*errors.AppErrors); ok && appErr.Code == 404 {
			if createdErr := s.repo.Create(post); createdErr != nil {
				return createdErr
			}
			return nil
		}
		return err
	}

	// Found an existing post — duplicate
	if existing != nil {
		return errors.Conflict("Post with the same title already exists", "Attempted to create a duplicate post")
	}

	return nil
}

func (s *postService) GetByID(id uint) (*models.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return post, nil

}

func (s *postService) GetAll(search, categoryID, authorID string, offset, limit int) ([]models.Post, int64, error) {
	posts, count, err := s.repo.FindAll(search, categoryID, authorID, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return posts, count, nil
}

func (s *postService) GetByAuthorID(authorID string, offset, limit int) ([]models.Post, int64, error) {
	posts, count, err := s.repo.FindAll("", "", authorID, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return posts, count, nil
}

func (s *postService) Update(post *models.Post, userID uint) error {
	if post.AuthorID != userID {
		return errors.Forbidden("You are not authorized to edit this post", "Tried to edit unauthorized post")
	}

	err := s.repo.Update(post)
	if err != nil {
		return err
	}
	return nil
}

func (s *postService) Delete(post *models.Post, userID uint) error {
	if post.AuthorID != userID {
		return errors.Forbidden("You are not authorized to delete this post", "Tried to delete unauthorized post")
	}
	err := s.repo.Delete(post)
	if err != nil {
		return err
	}
	return nil
}
