package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable             = "users"
	rolesTable             = "roles"
	sessionsTable          = "sessions"
	recipesTable           = "recipes"
	categoriesTable        = "categories"
	shoppingListTable      = "shopping_list"
	usersRecipesTable      = "users_recipes"
	recipesCategoriesTable = "recipes_categories"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	return db, nil
}
