package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/models"
	"time"
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

func (r *RecipesPostgres) GetRecipesByUser(userId int) ([]models.Recipe, error) {
	var recipes []models.Recipe
	query := fmt.Sprintf("SELECT %[1]v.*, %[2]v.favourite, %[2]v.liked FROM %[1]v LEFT JOIN %[2]v ON " +
		"%[2]v.recipe_id=%[1]v.recipe_id WHERE %[2]v.user_id=$1",
		recipesTable, usersRecipesTable)
	var ingredients []byte
	var cooking []byte
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return []models.Recipe{}, err
	}
	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Servings, &recipe.Time, &recipe.Calories,
			&ingredients, &cooking, &recipe.Preview, &recipe.Visibility, &recipe.Encrypted, &recipe.CreationTimestamp,
			&recipe.UpdateTimestamp, &recipe.Favourite, &recipe.Liked)
		if err != nil {
			return []models.Recipe{}, err
		}
		if err := json.Unmarshal(ingredients, &recipe.Ingredients); err != nil {
			return []models.Recipe{}, err
		}
		if err := json.Unmarshal(cooking, &recipe.Cooking); err != nil {
			return []models.Recipe{}, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
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

	var cookingTime int16 = 15
	if recipe.Servings >1  {
		cookingTime = recipe.Time
	}

	createRecipeQuery := fmt.Sprintf("INSERT INTO %s (name, owner_id, servings, cookingTime, calories, ingredients," +
		"cooking, preview, visibility, encrypted) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING recipe_id",
		recipesTable)
	row := tx.QueryRow(createRecipeQuery, recipe.Name, recipe.OwnerId, servings, cookingTime, recipe.Calories, ingredients,
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

func (r *RecipesPostgres) GetRecipeById(recipeId, userId int) (models.Recipe, error) {
	var recipe models.Recipe
	query := fmt.Sprintf("SELECT %[1]v.*, %[2]v.favourite, %[2]v.liked FROM %[1]v LEFT JOIN %[2]v ON " +
		"%[2]v.user_id=$1 AND %[1]v.recipe_id=%[2]v.recipe_id WHERE %[1]v.recipe_id=$2",
		recipesTable, usersRecipesTable)
	var ingredients []byte
	var cooking []byte
	row := r.db.QueryRow(query, userId, recipeId)
	if err := row.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Servings, &recipe.Time, &recipe.Calories, &ingredients,
		&cooking, &recipe.Preview, &recipe.Visibility, &recipe.Encrypted, &recipe.CreationTimestamp, &recipe.UpdateTimestamp, &recipe.Favourite, &recipe.Liked); err != nil {
		return models.Recipe{}, err
	}
	if err := json.Unmarshal(ingredients, &recipe.Ingredients); err != nil {
		return models.Recipe{}, err
	}
	if err := json.Unmarshal(cooking, &recipe.Cooking); err != nil {
		return models.Recipe{}, err
	}
	return recipe, nil
}

func (r *RecipesPostgres) GetRecipeOwnerId(recipeId int) (int, error) {
	var userId int
	query := fmt.Sprintf("SELECT owner_id FROM %s WHERE recipe_id=$1", recipesTable)
	err := r.db.Get(&userId, query, recipeId)
	return userId, err
}

func (r *RecipesPostgres) UpdateRecipe(recipe models.Recipe, userId int) error {
	ingredients, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		return err
	}

	cooking, err := json.Marshal(recipe.Cooking)
	if err != nil {
		return err
	}

	visibility := "private"
	if recipe.Visibility == r.visibilities[1] || recipe.Visibility == r.visibilities[2]  {
		visibility = recipe.Visibility
	}

	var servings int16 = 1
	if recipe.Servings >1  {
		servings = recipe.Servings
	}

	var cookingTime int16 = 15
	if recipe.Servings >1  {
		cookingTime = recipe.Time
	}

	query := fmt.Sprintf("UPDATE %s SET name=$1, servings=$2, time=$3, calories=$4, ingredients=$5, " +
		"cooking=$6, preview=$7, visibility=$8, encrypted=$9, update_timestamp=$10 WHERE recipe_id=$11 AND owner_id=$12",
		recipesTable)
	_, err = r.db.Exec(query, recipe.Name, servings, cookingTime, recipe.Calories, ingredients,
		cooking, recipe.Preview, visibility, recipe.Encrypted, time.Now(), recipe.Id, userId)
	return err
}

func (r *RecipesPostgres) DeleteRecipe(recipeId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE recipe_id=$1", recipesTable)
	_, err := r.db.Exec(query, recipeId)
	return err
}

func (r *RecipesPostgres) DeleteRecipeLink(recipeId, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE recipe_id=$1 AND user_id=$2", usersRecipesTable)
	_, err := r.db.Exec(query, recipeId, userId)
	return err
}