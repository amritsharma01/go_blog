package repositories

import (
	"crud_api/errors"
	"crud_api/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	FindByName(name string) (*models.Category, error)
	List(limit, offset int) ([]models.Category, int64, error)
	Delete(cat *models.Category) error
	FindByID(id uint) (*models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	if err := r.db.Create(category).Error; err != nil {
		return errors.Internal("Failed to create category", err)
	}
	return nil
}

func (r *categoryRepository) FindByName(name string) (*models.Category, error) {
	var category models.Category
	if err := r.db.Where("name = ?", name).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Category not found")
		}
		return nil, errors.Internal("Failed to find category by name", err)
	}
	return &category, nil
}

func (r *categoryRepository) List(limit, offset int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	if err := r.db.Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, errors.Internal("Failed to count categories", err)
	}

	if err := r.db.Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, errors.Internal("Failed to list categories", err)
	}

	return categories, total, nil
}

func (r *categoryRepository) Delete(cat *models.Category) error {
	if err := r.db.Delete(cat).Error; err != nil {
		return errors.Internal("Failed to delete category", err)
	}
	return nil
}

func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Category not found")
		}
		return nil, errors.Internal("Failed to fetch category by ID", err)
	}
	return &category, nil
}
