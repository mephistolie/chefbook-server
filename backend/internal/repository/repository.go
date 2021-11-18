package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository/postgres"
)

type Users interface {
	CreateUser(user models.AuthData, activationLink uuid.UUID) (int, error)
	GetUserById(id int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByCredentials(email, password string) (models.User, error)
	GetByRefreshToken(refreshToken string) (models.User, error)
	ActivateUser(activationLink uuid.UUID) error
	CreateSession(session models.Session) error
	UpdateSession(session models.Session, oldRefreshToken string) error
	DeleteSession(refreshToken string) error
	ChangePassword(user models.AuthData) error
}

type Recipes interface {
	GetRecipesByUser(userId int) ([]models.Recipe, error)
	GetRecipeOwnerId(recipeId int) (int, error)
	CreateRecipe(recipe models.Recipe) (int, error)
	GetRecipeById(recipeId int, userId int) (models.Recipe, error)
	UpdateRecipe(recipe models.Recipe, userId int) error
	DeleteRecipe(recipeId int) error
	DeleteRecipeLink(recipeId, userId int) error
	MarkRecipeFavourite(recipeId, userId int, isFavourite bool) error
}

type Categories interface {
	GetCategoriesByUser(userId int) ([]models.Category, error)
	AddCategory(category models.Category) (int, error)
	GetCategoryById(categoryId, userId int) (models.Category, error)
	UpdateCategory(category models.Category) error
	DeleteCategory(categoryId, userId int) error
}

type ShoppingList interface {
	GetShoppingList(userId int) ([]models.Purchase, error)
	SetShoppingList(shoppingList []models.Purchase, userId int) error
}

type Repository struct {
	Users
	Recipes
	Categories
	ShoppingList
}

func NewRepositories(db *sqlx.DB) *Repository {
	return &Repository{
		Users:        postgres.NewUsersPostgres(db),
		Recipes:      postgres.NewRecipesPostgres(db),
		Categories:   postgres.NewCategoriesPostgres(db),
		ShoppingList: postgres.NewShoppingListPostgres(db),
	}
}
