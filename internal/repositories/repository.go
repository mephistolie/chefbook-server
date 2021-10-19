package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/models"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(email, password string) (models.User, error)
}

type Recipes interface {

}

type Repository struct {
	Authorization
	Recipes
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}