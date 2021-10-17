package repositories

import "github.com/jmoiron/sqlx"

type Authorization interface {

}

type Recipes interface {

}

type Repository struct {
	Authorization
	Recipes
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}