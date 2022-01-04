package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/models"
	"time"
)

type RecipesPostgres struct {
	db *sqlx.DB
}

func NewRecipesPostgres(db *sqlx.DB) *RecipesPostgres {
	return &RecipesPostgres{
		db: db,
	}
}

func (r *RecipesPostgres) GetRecipesByUser(userId int) ([]models.Recipe, error) {
	var recipes []models.Recipe
	query := fmt.Sprintf("SELECT %[1]v.*, %[2]v.favourite, (SELECT EXISTS (SELECT 1 FROM %[3]v WHERE %[3]v.recipe_id=%[1]v.recipe_id AND user_id=$1)) as liked, %[4]v.username FROM %[1]v LEFT JOIN %[2]v ON " +
		"%[2]v.recipe_id=%[1]v.recipe_id LEFT JOIN %[4]v ON %[4]v.user_id=%[1]v.owner_id WHERE %[2]v.user_id=$1",
		recipesTable, usersRecipesTable, likesTable, usersTable)
	var ingredients []byte
	var cooking []byte
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return []models.Recipe{}, err
	}
	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Description, &recipe.Likes, &recipe.Servings, &recipe.Time, &recipe.Calories,
			&ingredients, &cooking, &recipe.Preview, &recipe.Visibility, &recipe.Encrypted, &recipe.CreationTimestamp,
			&recipe.UpdateTimestamp, &recipe.Favourite, &recipe.Liked, &recipe.OwnerName)
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

	createRecipeQuery := fmt.Sprintf("INSERT INTO %s (name, owner_id, description, servings, time, calories, ingredients," +
		"cooking, preview, visibility, encrypted) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING recipe_id",
		recipesTable)
	row := tx.QueryRow(createRecipeQuery, recipe.Name, recipe.OwnerId, recipe.Description, recipe.Servings, recipe.Time, recipe.Calories, recipe.Ingredients,
		recipe.Cooking, recipe.Preview, recipe.Visibility, recipe.Encrypted)
	if err := row.Scan(&id); err != nil {
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, err
	}

	createRecipeLinkQuery := fmt.Sprintf("INSERT INTO %s (user_id, recipe_id) values ($1, $2) RETURNING user_id",
		usersRecipesTable)
	if _, err := tx.Exec(createRecipeLinkQuery, recipe.OwnerId, id); err != nil {
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, err
	}
	err = tx.Commit()
	return id, nil
}

func (r *RecipesPostgres) SetRecipeCategories(categoriesIds []int, recipeId, userId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	clearCategoriesQuery := fmt.Sprintf("DELETE FROM %s WHERE recipe_id=$1 AND user_id=$2", recipesCategoriesTable)
	if _, err := tx.Exec(clearCategoriesQuery, recipeId, userId); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if len(categoriesIds) > 0 {
		categoriesArrayString := "("
		for _, categoryId := range categoriesIds {
			categoriesArrayString += fmt.Sprintf("%d, ", categoryId)
		}
		categoriesArrayString = categoriesArrayString[:len(categoriesArrayString)-2] + ")"
		addCategoriesQuery := fmt.Sprintf("INSERT INTO %[1]v (recipe_id, category_id, user_id) " +
			"SELECT %[2]v.recipe_id, %[3]v.category_id, %[3]v.user_id FROM %[3]v LEFT JOIN %[2]v ON %[2]v.recipe_id=$1 WHERE category_id IN %[4]v AND user_id=$2",
			recipesCategoriesTable, recipesTable, categoriesTable, categoriesArrayString)
		if _, err := tx.Exec(addCategoriesQuery, recipeId, userId); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *RecipesPostgres) GetRecipeById(recipeId, userId int) (models.Recipe, error) {
	var recipe models.Recipe
	query := fmt.Sprintf("SELECT %[1]v.*, %[2]v.favourite, (SELECT EXISTS (SELECT 1 FROM %[3]v WHERE %[3]v.recipe_id=%[1]v.recipe_id AND user_id=1)) as liked, %[4]v.username FROM %[1]v LEFT JOIN %[2]v ON " +
		"%[2]v.user_id=$1 AND %[1]v.recipe_id=%[2]v.recipe_id LEFT JOIN users ON %[4]v.user_id=%[1]v.owner_id WHERE %[1]v.recipe_id=$2",
		recipesTable, usersRecipesTable, likesTable, usersTable)
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
	query := fmt.Sprintf("UPDATE %s SET name=$1, description=$2, servings=$3, time=$4, calories=$5, ingredients=$6, " +
		"cooking=$7, preview=$8, visibility=$9, encrypted=$10, update_timestamp=$11 WHERE recipe_id=$12 AND owner_id=$13",
		recipesTable)
	_, err := r.db.Exec(query, recipe.Name, recipe.Description, recipe.Servings, recipe.Time, recipe.Calories, recipe.Ingredients,
		recipe.Cooking, recipe.Preview, recipe.Visibility, recipe.Encrypted, time.Now(), recipe.Id, userId)
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

func (r *RecipesPostgres) SetRecipeLike(recipeId, userId int, isLiked bool) error {
	var exists bool
	checkLikeQuery := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE recipe_id=$1 AND user_id=$2)", likesTable)

	err := r.db.QueryRow(checkLikeQuery, recipeId, userId).Scan(&exists); if err != nil && err != sql.ErrNoRows {
		return err
	}
	if (exists && isLiked) || (!exists && !isLiked) {
		return models.ErrRecipeLikeSetAlready
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var likeCountQuery string
	if isLiked {
		likeRecipeQuery := fmt.Sprintf("INSERT INTO %s (recipe_id, user_id) values ($1, $2)", likesTable)
		if _, err := tx.Exec(likeRecipeQuery, recipeId, userId); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
		likeCountQuery = fmt.Sprintf("UPDATE %s SET likes=likes+1 WHERE recipe_id=$1", recipesTable)
	} else {
		unlikeRecipeQuery := fmt.Sprintf("DELETE FROM %s WHERE recipe_id=$1 AND user_id=$2", likesTable)
		if _, err := tx.Exec(unlikeRecipeQuery, recipeId, userId); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
		likeCountQuery = fmt.Sprintf("UPDATE %s SET likes=likes-1 WHERE recipe_id=$1", recipesTable)
	}

	if _, err := tx.Exec(likeCountQuery, recipeId); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (r *RecipesPostgres) MarkRecipeFavourite(recipeId, userId int, isFavourite bool) error {
	query := fmt.Sprintf("UPDATE %s SET favourite=$1 WHERE recipe_id=$2 AND user_id=$3", usersRecipesTable)
	_, err := r.db.Exec(query, isFavourite, recipeId, userId)
	return err
}