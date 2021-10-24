package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/models"
)

type RecipesPostgres struct {
	db *sqlx.DB
	visibilities [3]string
}

func NewRecipesPostgres(db *sqlx.DB) *RecipesPostgres {
	return &RecipesPostgres{
		db: db,
		visibilities: [3]string{"private", "shared", "public"},
	}
}

func (r *RecipesPostgres) CreateRecipe(recipe models.Recipe) (int, error) {
	var id int
	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}

	ingredients, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		return -1, err
	}

	cooking, err := json.Marshal(recipe.Cooking)
	if err != nil {
		return -1, err
	}

	visibility := "private"
	if recipe.Visibility == r.visibilities[1] || recipe.Visibility == r.visibilities[2]  {
		visibility = recipe.Visibility
	}

	var servings int16 = 1
	if recipe.Servings >1  {
		servings = recipe.Servings
	}

	var time int16 = 15
	if recipe.Servings >1  {
		time = recipe.Time
	}

	createRecipeQuery := fmt.Sprintf("INSERT INTO %s (name, owner_id, servings, time, calories, ingredients," +
		"cooking, preview, visibility, encrypted) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING recipe_id",
		recipesTable)
	row := tx.QueryRow(createRecipeQuery, recipe.Name, recipe.OwnerId, servings, time, recipe.Calories, ingredients,
		cooking, recipe.Preview, visibility, recipe.Encrypted)
	if err := row.Scan(&id); err != nil {
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, err
	}

	createRoleQuery := fmt.Sprintf("INSERT INTO %s (user_id, recipe_id) values ($1, $2) RETURNING user_id",
		usersRecipesTable)
	if _, err := tx.Exec(createRoleQuery, recipe.OwnerId, id); err != nil {
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, err
	}

	err = tx.Commit()
	return id, nil
}