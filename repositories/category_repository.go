package repositories

import (
	"crud_api/errors"
	"crud_api/models"
	"fmt"

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
		return errors.Internal(
			"Unable to create category",
			"Database error while creating category",
			err,
		)
	}
	return nil
}

func (r *categoryRepository) FindByName(name string) (*models.Category, error) {
	var category models.Category
	if err := r.db.Where("name = ?", name).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(
				"Category not found",
				fmt.Sprintf("Category with name '%s' not found", name),
			)
		}
		return nil, errors.Internal(
			"Unable to search for category",
			"Database error while searching for category by name",
			err,
		)
	}
	return &category, nil
}

func (r *categoryRepository) List(limit, offset int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	if err := r.db.Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, errors.Internal(
			"Unable to retrieve categories",
			"Database error while counting categories",
			err,
		)
	}

	if err := r.db.Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, errors.Internal(
			"Unable to retrieve categories",
			"Database error while listing categories",
			err,
		)
	}

	return categories, total, nil
}

func (r *categoryRepository) Delete(cat *models.Category) error {
	if err := r.db.Delete(cat).Error; err != nil {
		return errors.Internal(
			"Unable to delete category",
			"Database error while deleting category",
			err,
		)
	}
	return nil
}

func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound(
				"Category not found",
				fmt.Sprintf("Category with ID %d not found", id),
			)
		}
		return nil, errors.Internal(
			"Unable to retrieve category",
			"Database error while finding category by ID",
			err,
		)
	}
	return &category, nil
}
