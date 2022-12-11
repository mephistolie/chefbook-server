package repository

import "github.com/mephistolie/chefbook-server/internal/entity"

type Category interface {
	GetUserCategories(userId string) []entity.Category
	GetRecipeCategories(recipeId, userId string) []entity.Category
	GetCategory(categoryId string) (entity.Category, error)
	GetCategoryOwnerId(categoryId string) (string, error)
	CreateCategory(category entity.CategoryInput, userId string) (string, error)
	UpdateCategory(categoryId string, category entity.CategoryInput) error
	DeleteCategory(categoryId string) error
}
