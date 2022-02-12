package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/model"
)

type RecipeSharingPostgres struct {
	db *sqlx.DB
}

func NewRecipeSharingPostgres(db *sqlx.DB) *RecipeSharingPostgres {
	return &RecipeSharingPostgres{
		db: db,
	}
}

func (r *RecipeSharingPostgres) GetRecipeUserList(recipeId int) ([]model.UserInfo, error) {
	query := fmt.Sprintf("SELECT %[1]v.user_id, %[2]v.username, %[2]v.avatar, %[2]v.premium "+
		"LEFT JOIN %[2]v ON %[2]v.user_id=%[1]v.user_id WHERE %[1]v.user_id=$1", usersRecipesTable, usersTable)
	rows, err := r.db.Query(query, recipeId)
	if err != nil {
		return []model.UserInfo{}, err
	}
	var users []model.UserInfo
	for rows.Next() {
		var user model.UserInfo
		err := rows.Scan(&user.Id, &user.Username, &user.Avatar, &user.Premium)
		if err != nil {
			return []model.UserInfo{}, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *RecipeSharingPostgres) SetUserPublicKeyForRecipe(recipeId int, userId int, userKey string) error {
	query := fmt.Sprintf("UPDATE %s SET user_key=$1 WHERE recipe_id=$2 AND user_id=$3", usersRecipesTable)
	_, err := r.db.Exec(query, userKey, recipeId, userId)
	return err
}

func (r *RecipeSharingPostgres) SetUserPrivateKeyForRecipe(recipeId int, userId int, recipeKey string) error {
	query := fmt.Sprintf("UPDATE %s SET recipe_key=$1 WHERE recipe_id=$2 AND user_id=$3", usersRecipesTable)
	_, err := r.db.Exec(query, recipeKey, recipeId, userId)
	return err
}