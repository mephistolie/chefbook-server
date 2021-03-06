package service

import (
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/service/interface/repository"
)

type CategoriesService struct {
	repo repository.Category
}

func NewCategoriesService(repo repository.Category) *CategoriesService {
	return &CategoriesService{
		repo: repo,
	}
}

func (s *CategoriesService) GetUserCategories(userId int) []entity.Category {
	return s.repo.GetUserCategories(userId)
}

func (s *CategoriesService) GetRecipeCategories(recipeId, userId int) []entity.Category {
	return s.repo.GetRecipeCategories(recipeId, userId)
}

func (s *CategoriesService) CreateCategory(category entity.CategoryInput, userId int) (int, error) {
	return s.repo.CreateCategory(category, userId)
}

func (s *CategoriesService) GetCategory(categoryId int, userId int) (entity.Category, error) {
	category, err := s.repo.GetCategory(categoryId)
	if err != nil {
		return entity.Category{}, err
	}

	if category.UserId != userId {
		return entity.Category{}, failure.AccessDenied
	}

	return category, nil
}

func (s *CategoriesService) UpdateCategory(categoryId int, category entity.CategoryInput, userId int) error {
	ownerId, err := s.repo.GetCategoryOwnerId(categoryId)
	if err != nil {
		return err
	}
	if ownerId != userId {
		return failure.AccessDenied
	}

	return s.repo.UpdateCategory(categoryId, category)
}

func (s *CategoriesService) DeleteCategory(categoryId, userId int) error {
	category, err := s.repo.GetCategory(categoryId)
	if err != nil {
		return err
	}

	if category.UserId != userId {
		return failure.AccessDenied
	}

	return s.repo.DeleteCategory(categoryId)
}
