package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	models2 "github.com/mephistolie/chefbook-server/internal/models"
	postgres2 "github.com/mephistolie/chefbook-server/internal/repository/postgres"
)

type Users interface {
	CreateUser(user models2.AuthData, activationLink uuid.UUID) (int, error)
	GetUserById(id int) (models2.User, error)
	GetUserByEmail(email string) (models2.User, error)
	GetUserByCredentials(email, password string) (models2.User, error)
	GetByRefreshToken(refreshToken string) (models2.User, error)
	ActivateUser(activationLink uuid.UUID) error
	CreateSession(session models2.Session) error
	UpdateSession(session models2.Session, oldRefreshToken string) error
	DeleteSession(refreshToken string) error
	ChangePassword(user models2.AuthData) error
}

type Recipes interface {
	GetRecipesByUser(userId int) ([]models2.Recipe, error)
	GetRecipeOwnerId(recipeId int) (int, error)
	CreateRecipe(recipe models2.Recipe) (int, error)
	GetRecipeById(recipeId int, userId int) (models2.Recipe, error)
	UpdateRecipe(recipe models2.Recipe, userId int) error
	DeleteRecipe(recipeId int) error
	DeleteRecipeLink(recipeId, userId int) error
}

type Repository struct {
	Users
	Recipes
}

func NewRepositories(db *sqlx.DB) *Repository {
	return &Repository{
		Users:   postgres2.NewUsersPostgres(db),
		Recipes: postgres2.NewRecipesPostgres(db),
	}
}