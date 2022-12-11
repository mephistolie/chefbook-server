package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/repository/postgres/dto"
	"time"
)

type RecipeOwnershipPostgres struct {
	db *sqlx.DB
}

func NewRecipeOwnershipPostgres(db *sqlx.DB) *RecipeOwnershipPostgres {
	return &RecipeOwnershipPostgres{
		db: db,
	}
}

func (r *RecipeOwnershipPostgres) CreateRecipe(recipe entity.RecipeInput, userId uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	tx, err := r.db.Begin()
	if err != nil {
		logRepoError(err)
		return uuid.UUID{}, failure.Unknown
	}

	bsonIngredients, err := json.Marshal(dto.NewIngredients(recipe.Ingredients))
	if err != nil {
		logRepoError(err)
		return uuid.UUID{}, failure.Unknown
	}
	bsonCooking, err := json.Marshal(dto.NewCooking(recipe.Cooking))
	if err != nil {
		logRepoError(err)
		return uuid.UUID{}, failure.Unknown
	}

	createRecipeQuery := fmt.Sprintf(`
			INSERT INTO %s
				(name, owner_id, language, description, servings, time, calories, protein, fats, carbohydrates, ingredients,
				cooking, preview, visibility, encrypted)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			RETURNING
				recipe_id
		`, recipesTable)

	row := tx.QueryRow(createRecipeQuery, recipe.Name, userId, recipe.Language, recipe.Description, recipe.Servings,
		recipe.Time, recipe.Calories, recipe.Macronutrients.Protein, recipe.Macronutrients.Fats, recipe.Macronutrients.Carbohydrates,
		bsonIngredients, bsonCooking, recipe.Preview, recipe.Visibility, recipe.IsEncrypted)
	if err := row.Scan(&id); err != nil {
		logRepoError(err)
		if err := tx.Rollback(); err != nil {
			logRepoError(err)
			return uuid.UUID{}, failure.Unknown
		}
		return uuid.UUID{}, failure.UnableCreateRecipe
	}

	createRecipeLinkQuery := fmt.Sprintf(`
			INSERT INTO %s (user_id, recipe_id)
			VALUES ($1, $2)
			RETURNING user_id
		`, usersRecipesTable)

	if _, err := tx.Exec(createRecipeLinkQuery, userId, id); err != nil {
		logRepoError(err)
		if err := tx.Rollback(); err != nil {
			logRepoError(err)
			return uuid.UUID{}, err
		}
		return uuid.UUID{}, failure.UnableCreateRecipe
	}

	if err = tx.Commit(); err != nil {
		logRepoError(err)
		return uuid.UUID{}, failure.UnableCreateRecipe
	}

	return id, nil
}

func (r *RecipeOwnershipPostgres) UpdateRecipe(recipeId uuid.UUID, recipe entity.RecipeInput) error {

	bsonIngredients, err := json.Marshal(dto.NewIngredients(recipe.Ingredients))
	if err != nil {
		logRepoError(err)
		return failure.Unknown
	}
	bsonCooking, err := json.Marshal(dto.NewCooking(recipe.Cooking))
	if err != nil {
		logRepoError(err)
		return failure.Unknown
	}

	updateRecipeQuery := fmt.Sprintf(`
			UPDATE
				%s
			SET
				name=$1, language=$2, description=$3, servings=$4, time=$5, calories=$6, protein=$7, fats=$8,
				carbohydrates=$9, ingredients=$10, cooking=$11, preview=$12, visibility=$13, encrypted=$14, update_timestamp=$15
			WHERE
				recipe_id=$16
		`, recipesTable)

	if _, err := r.db.Exec(updateRecipeQuery, recipe.Name, recipe.Language, recipe.Description, recipe.Servings,
		recipe.Time, recipe.Calories, recipe.Macronutrients.Protein,
		recipe.Macronutrients.Fats, recipe.Macronutrients.Carbohydrates, bsonIngredients, bsonCooking, recipe.Preview,
		recipe.Visibility, recipe.IsEncrypted, time.Now().UTC(), recipeId); err != nil {
		logRepoError(err)
		return failure.RecipeNotFound
	}

	return nil
}

func (r *RecipeOwnershipPostgres) DeleteRecipe(recipeId uuid.UUID) error {

	deleteRecipeQuery := fmt.Sprintf(`
			DELETE FROM %s
			WHERE recipe_id=$1
		`, recipesTable)

	if _, err := r.db.Exec(deleteRecipeQuery, recipeId); err != nil {
		logRepoError(err)
		return failure.RecipeNotFound
	}

	return nil
}
