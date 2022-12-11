package service

import "github.com/mephistolie/chefbook-server/internal/entity"

type Category interface {
	GetUserCategories(userId string) []entity.Category
	GetRecipeCategories(recipeId, userId string) []entity.Category
	CreateCategory(category entity.CategoryInput, userId string) (string, error)
	GetCategory(categoryId string, userId string) (entity.Category, error)
	UpdateCategory(categoryId string, category entity.CategoryInput, userId string) error
	DeleteCategory(categoryId, userId string) error
}
