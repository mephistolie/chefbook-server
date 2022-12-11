package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Category interface {
	GetUserCategories(userId uuid.UUID) []entity.Category
	GetRecipeCategories(recipeId, userId uuid.UUID) []entity.Category
	CreateCategory(category entity.CategoryInput, userId uuid.UUID) (uuid.UUID, error)
	GetCategory(categoryId uuid.UUID, userId uuid.UUID) (entity.Category, error)
	UpdateCategory(categoryId uuid.UUID, category entity.CategoryInput, userId uuid.UUID) error
	DeleteCategory(categoryId, userId uuid.UUID) error
}
