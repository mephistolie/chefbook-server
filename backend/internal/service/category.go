package service

import (
	"github.com/google/uuid"
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

func (s *CategoriesService) GetUserCategories(userId uuid.UUID) []entity.Category {
	return s.repo.GetUserCategories(userId)
}

func (s *CategoriesService) GetRecipeCategories(recipeId, userId uuid.UUID) []entity.Category {
	return s.repo.GetRecipeCategories(recipeId, userId)
}

func (s *CategoriesService) CreateCategory(category entity.CategoryInput, userId uuid.UUID) (uuid.UUID, error) {
	return s.repo.CreateCategory(category, userId)
}

func (s *CategoriesService) GetCategory(categoryId, userId uuid.UUID) (entity.Category, error) {
	category, err := s.repo.GetCategory(categoryId)
	if err != nil {
		return entity.Category{}, err
	}

	if category.UserId != userId {
		return entity.Category{}, failure.AccessDenied
	}

	return category, nil
}

func (s *CategoriesService) UpdateCategory(categoryId uuid.UUID, category entity.CategoryInput, userId uuid.UUID) error {
	ownerId, err := s.repo.GetCategoryOwnerId(categoryId)
	if err != nil {
		return err
	}
	if ownerId != userId {
		return failure.AccessDenied
	}

	return s.repo.UpdateCategory(categoryId, category)
}

func (s *CategoriesService) DeleteCategory(categoryId, userId uuid.UUID) error {
	category, err := s.repo.GetCategory(categoryId)
	if err != nil {
		return err
	}

	if category.UserId != userId {
		return failure.AccessDenied
	}

	return s.repo.DeleteCategory(categoryId)
}
