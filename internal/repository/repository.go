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
	CreateRecipe(recipe models.Recipe) (int, error)
}

type Repository struct {
	Users
	Recipes
}

func NewRepositories(db *sqlx.DB) *Repository {
	return &Repository{
		Users: postgres.NewUsersPostgres(db),
		Recipes: postgres.NewRecipesPostgres(db),
	}
}