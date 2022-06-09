package service

import "github.com/mephistolie/chefbook-server/internal/entity"

type Category interface {
	GetUserCategories(userId int) []entity.Category
	GetRecipeCategories(recipeId, userId int) []entity.Category
	CreateCategory(category entity.CategoryInput, userId int) (int, error)
	GetCategory(categoryId int, userId int) (entity.Category, error)
	UpdateCategory(categoryId int, category entity.CategoryInput, userId int) error
	DeleteCategory(categoryId, userId int) error
}
