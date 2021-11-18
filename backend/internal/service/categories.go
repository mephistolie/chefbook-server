package service

import (
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository"
)

type CategoriesService struct {
	repo repository.Categories
}

func NewCategoriesService(repo repository.Categories) *CategoriesService {
	return &CategoriesService{
		repo: repo,
	}
}

func (s *CategoriesService) GetCategoriesByUser(userId int) ([]models.Category, error) {
	categories, err := s.repo.GetCategoriesByUser(userId)
	if categories == nil {
		categories = []models.Category{}
	}
	return categories, err
}

func (s *CategoriesService) AddCategory(category models.Category) (int, error) {
	return s.repo.AddCategory(category)
}

func (s *CategoriesService) GetCategoryById(categoryId, userId int) (models.Category, error) {
	category, err := s.repo.GetCategoryById(categoryId, userId)
	if err != nil {
		return models.Category{}, models.ErrCategoryNotFound
	}
	return category, err
}

func (s *CategoriesService) UpdateCategory(category models.Category) error {
	return s.repo.UpdateCategory(category)
}

func (s *CategoriesService) DeleteCategory(recipeId, userId int) error {
	return s.repo.DeleteCategory(recipeId, userId)
}