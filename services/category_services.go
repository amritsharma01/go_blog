package services

import (
	"crud_api/models"
	"crud_api/repositories"
)

type CategoryService interface {
	AddCategory(category *models.Category) error
	GetCategories(limit, offset int) ([]models.Category, int64, error)
	CategoryExists(name string) (bool, error)
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
	return s.repo.Create(category)
}

func (s *categoryService) GetCategories(limit, offset int) ([]models.Category, int64, error) {
	return s.repo.List(limit, offset)
}

func (s *categoryService) CategoryExists(name string) (bool, error) {
	cat, err := s.repo.FindByName(name)
	if err != nil {
		return false, nil // not found is not an error for existence check
	}
	return cat != nil, nil
}

func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	return s.repo.FindByID(id)
}

func (s *categoryService) DeleteCategory(category *models.Category) error {
	return s.repo.Delete(category)
}
