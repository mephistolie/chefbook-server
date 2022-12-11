package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Category interface {
	GetUserCategories(userId uuid.UUID) []entity.Category
	GetRecipeCategories(recipeId, userId uuid.UUID) []entity.Category
	GetCategory(categoryId uuid.UUID) (entity.Category, error)
	GetCategoryOwnerId(categoryId uuid.UUID) (uuid.UUID, error)
	CreateCategory(category entity.CategoryInput, userId uuid.UUID) (uuid.UUID, error)
	UpdateCategory(categoryId uuid.UUID, category entity.CategoryInput) error
	DeleteCategory(categoryId uuid.UUID) error
}
