package repository

import "github.com/mephistolie/chefbook-server/internal/entity"

type Category interface {
	GetUserCategories(userId int) []entity.Category
	GetRecipeCategories(recipeId, userId int) []entity.Category
	GetCategory(categoryId int) (entity.Category, error)
	GetCategoryOwnerId(categoryId int) (int, error)
	CreateCategory(category entity.CategoryInput, userId int) (int, error)
	UpdateCategory(categoryId int, category entity.CategoryInput) error
	DeleteCategory(categoryId int) error
}
