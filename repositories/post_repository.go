package repositories

import (
	"crud_api/errors"
	"crud_api/models"
	"log"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	FindByID(id uint) (*models.Post, error)
	FindAll(search string, categoryID, authorID string, offset, limit int) ([]models.Post, int64, error)
	FindByAuthorID(authorID uint, offset, limit int) ([]models.Post, int64, error)
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
		return errors.New(500, "Failed to create post", err)
	}
	return nil
}

func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Preload("Author").Preload("Category").First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewWithMessage(404, "Post not found")
		}
		return nil, errors.New(500, "Failed to fetch post by ID", err)
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
		return nil, 0, errors.New(500, "Failed to count posts", err)
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, 0, errors.New(500, "Failed to fetch posts", err)
	}

	return posts, count, nil
}

func (r *postRepository) FindByAuthorID(authorID uint, offset, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var count int64

	query := r.db.Model(&models.Post{}).Where("author_id = ?", authorID).Preload("Author").Preload("Category")

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, errors.New(500, "Failed to count posts by author", err)
	}

	if err := query.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, 0, errors.New(500, "Failed to fetch posts by author", err)
	}

	return posts, count, nil
}

func (r *postRepository) Update(post *models.Post) error {
	log.Printf("Updating Post: Title=%s, Description=%s, CategoryID=%d", post.Title, post.Description, post.CategoryID)
	if err := r.db.Save(post).Error; err != nil {
		return errors.New(500, "Failed to update post", err)
	}
	return nil
}

func (r *postRepository) Delete(post *models.Post) error {
	if err := r.db.Delete(post).Error; err != nil {
		return errors.New(500, "Failed to delete post", err)
	}
	return nil
}

func (r *postRepository) FindDuplicate(title string, authorID uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Where("title = ? AND author_id = ?", title, authorID).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No duplicate is not an error
		}
		return nil, errors.New(500, "Failed to check for duplicate post", err)
	}
	return &post, nil
}
