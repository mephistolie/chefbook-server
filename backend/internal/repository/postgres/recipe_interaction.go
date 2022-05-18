package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type RecipeInteractionPostgres struct {
	db *sqlx.DB
}

func NewRecipeInteractionPostgres(db *sqlx.DB) *RecipeInteractionPostgres {
	return &RecipeInteractionPostgres{
		db: db,
	}
}

func (r *RecipeInteractionPostgres) SetRecipeCategories(categoriesIds []int, recipeId, userId int) error {
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
		addCategoriesQuery := fmt.Sprintf("INSERT INTO %[1]v (recipe_id, category_id, user_id) "+
			"SELECT %[2]v.recipe_id, %[3]v.category_id, %[3]v.user_id FROM %[3]v LEFT JOIN %[2]v ON %[2]v.recipe_id=$1 "+
			"WHERE category_id IN %[4]v AND user_id=$2",
			recipesCategoriesTable, recipesTable, categoriesTable, categoriesArrayString)
		if _, err := tx.Exec(addCategoriesQuery, recipeId, userId); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	updateUserTimestamp := fmt.Sprintf("UPDATE %s SET update_timestamp=$1 WHERE recipe_id=$2 AND user_id=$3", usersRecipesTable)
	if _, err := tx.Exec(updateUserTimestamp, time.Now(), recipeId, userId); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (r *RecipeInteractionPostgres) SetRecipeFavourite(recipeId, userId int, isFavourite bool) error {
	query := fmt.Sprintf("UPDATE %s SET favourite=$1, update_timestamp=$2 WHERE recipe_id=$3 AND user_id=$4", usersRecipesTable)
	_, err := r.db.Exec(query, isFavourite, time.Now(), recipeId, userId)
	return err
}

func (r *RecipeInteractionPostgres) SetRecipeLiked(recipeId, userId int, isLiked bool) error {
	var exists bool
	checkLikeQuery := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE recipe_id=$1 AND user_id=$2)", likesTable)

	err := r.db.QueryRow(checkLikeQuery, recipeId, userId).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if (exists && isLiked) || (!exists && !isLiked) {
		return nil
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