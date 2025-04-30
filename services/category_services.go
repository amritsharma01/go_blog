package services

import (
	"crud_api/errors"
	"crud_api/models"
	"crud_api/repositories"
)

type CategoryService interface {
	AddCategory(category *models.Category) error
	GetCategories(limit, offset int) ([]models.Category, int64, error)
	DeleteCategory(cat *models.Category) error
	GetByID(id uint) (*models.Category, error)
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) AddCategory(category *models.Category) error {
	existingCategory, err := s.repo.FindByName(category.Name)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == 404 {
			// Category does not exist, proceed to create it
		} else {
			return err
		}
	} else if existingCategory != nil {
		return errors.NewWithMessage(409, "Category already exists")
	}

	if err := s.repo.Create(category); err != nil {
		return errors.Internal("Failed to create category", err)
	}

	return nil
}

func (s *categoryService) GetCategories(limit, offset int) ([]models.Category, int64, error) {
	categories, total, err := s.repo.List(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) DeleteCategory(category *models.Category) error {
	if err := s.repo.Delete(category); err != nil {
		return errors.Internal("Failed to delete category", err)
	}
	return nil
}
