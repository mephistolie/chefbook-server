package service

import (
	"github.com/mephistolie/chefbook-server/internal/model"
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

func (s *CategoriesService) GetUserCategories(userId int) ([]model.Category, error) {
	categories, err := s.repo.GetUserCategories(userId)
	if categories == nil {
		categories = []model.Category{}
	}
	return categories, err
}

func (s *CategoriesService) GetRecipeCategories(recipeId, userId int) ([]model.Category, error) {
	return s.repo.GetRecipeCategories(recipeId, userId)
}

func (s *CategoriesService) AddCategory(category model.Category) (int, error) {
	return s.repo.AddCategory(category)
}

func (s *CategoriesService) GetCategoryById(categoryId, userId int) (model.Category, error) {
	category, err := s.repo.GetCategoryById(categoryId)
	if err != nil {
		return model.Category{}, model.ErrCategoryNotFound
	}
	if category.UserId != userId {
		return model.Category{}, model.ErrAccessDenied
	}
	return category, err
}

func (s *CategoriesService) UpdateCategory(category model.Category) error {
	return s.repo.UpdateCategory(category)
}

func (s *CategoriesService) DeleteCategory(categoryId, userId int) error {
	return s.repo.DeleteCategory(categoryId, userId)
}