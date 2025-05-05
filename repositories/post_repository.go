package repositories

import (
	"crud_api/errors"
	"crud_api/models"
	"fmt"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	FindByID(id uint) (*models.Post, error)
	FindAll(search, categoryID, authorID string, offset, limit int) ([]models.Post, int64, error)
	Update(post *models.Post) error
	Delete(post *models.Post) error
	FindDuplicate(title string, authorID uint) (*models.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db}
}

func (r *postRepository) Create(post *models.Post) error {
	if err := r.db.Create(post).Error; err != nil {
		return errors.Internal(
			"Unable to create post",
			"Database error while creating post",
			err,
		)
	}
	return nil
}

func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Preload("Author").Preload("Category").First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Post not found",
				fmt.Sprintf("Post with id '%d' not found", id),
			)
		}
		return nil, errors.Internal(
			"Unable to find post",
			"Database error while searching for post by ID",
			err)
	}
	return &post, nil
}

func (r *postRepository) FindAll(search, categoryID, authorID string, offset, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var count int64

	query := r.db.Model(&models.Post{}).Preload("Author").Preload("Category")

	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if authorID != "" {
		query = query.Where("author_id = ?", authorID)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, errors.Internal("unable to count posts",
			"Database error while counting posts", err)
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, 0, errors.Internal("Unable to retrieve posts", "Database error while retrieving posts", err)
	}

	return posts, count, nil
}

func (r *postRepository) Update(post *models.Post) error {

	if err := r.db.Save(post).Error; err != nil {
		return errors.Internal("unable to update post", "Database error while updating post", err)
	}
	return nil
}

func (r *postRepository) Delete(post *models.Post) error {
	if err := r.db.Delete(post).Error; err != nil {
		return errors.Internal("unable to delete post", "Database error while deleting post", err)
	}
	return nil
}

func (r *postRepository) FindDuplicate(title string, authorID uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Where("title = ? AND author_id = ?", title, authorID).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(
				"Post not found",
				fmt.Sprintf("Post with title '%s' not found", title),
			) // No duplicate found is not an error
		}
		return nil, errors.Internal(
			"Unable to check for duplicate post",
			"Database error while checking for duplicate post",
			err)
	}
	return &post, nil
}
