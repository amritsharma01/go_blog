package repositories

import (
	"crud_api/models"

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
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Preload("Author").Preload("Category").First(&post, id).Error; err != nil {
		return nil, err
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
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (r *postRepository) FindByAuthorID(authorID uint, offset, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var count int64

	query := r.db.Model(&models.Post{}).Where("author_id = ?", authorID).Preload("Author").Preload("Category")

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(post *models.Post) error {
	return r.db.Delete(post).Error
}

func (r *postRepository) FindDuplicate(title string, authorID uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Where("title = ? AND author_id = ?", title, authorID).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}
