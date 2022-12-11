package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
)

type RecipeSharingPostgres struct {
	db *sqlx.DB
}

func NewRecipeSharingPostgres(db *sqlx.DB) *RecipeSharingPostgres {
	return &RecipeSharingPostgres{
		db: db,
	}
}

func (r *RecipeSharingPostgres) GetRecipeUserList(recipeId string) ([]entity.ProfileInfo, error) {

	query := fmt.Sprintf(`
			SELECT
				%[1]v.user_id, %[2]v.username, %[2]v.registered, %[2]v.avatar, %[2]v.premium, %[2]v.broccoins
			FROM
				%[1]v
			LEFT JOIN
				%[2]v ON %[2]v.user_id=%[1]v.user_id
			WHERE
				%[1]v.recipe_id=$1
		`, usersRecipesTable, usersTable)

	rows, err := r.db.Query(query, recipeId)
	if err != nil {
		logRepoError(err)
		return []entity.ProfileInfo{}, failure.RecipeNotFound
	}

	var users []entity.ProfileInfo
	for rows.Next() {
		var user entity.ProfileInfo
		err := rows.Scan(&user.Id, &user.Username, &user.CreationTimestamp, &user.Avatar, &user.PremiumEndDate, &user.Broccoins)
		if err != nil {
			logRepoError(err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *RecipeSharingPostgres) GetUserPublicKey(recipeId, userId string) (string, error) {
	var key *string

	setUserPublicKeyQuery := fmt.Sprintf(`
			SELECT user_key
			FROM
				%[1]v
			WHERE recipe_id=$1 AND user_id=$2
		`, usersRecipesTable)

	if err := r.db.Get(&key, setUserPublicKeyQuery, recipeId, userId); err != nil || key == nil {
		if err != nil {
			logRepoError(err)
		}
		return "", failure.NoKey
	}

	return *key, nil
}

func (r *RecipeSharingPostgres) SetUserPublicKeyLink(recipeId string, userId string, userKey *string) error {

	setUserPublicKeyQuery := fmt.Sprintf(`
			UPDATE %s
			SET user_key=$1
			WHERE recipe_id=$2 AND user_id=$3
		`, usersRecipesTable)

	res, err := r.db.Exec(setUserPublicKeyQuery, userKey, recipeId, userId)
	if err != nil {
		logRepoError(err)
		return failure.RecipeNotInRecipeBook
	}

	if rowsCount, err := res.RowsAffected(); err == nil && rowsCount == 0 {
		return failure.RecipeNotInRecipeBook
	}

	return nil
}

func (r *RecipeSharingPostgres) GetUserRecipeKey(recipeId, userId string) (string, error) {
	var key *string

	setUserPublicKeyQuery := fmt.Sprintf(`
			SELECT recipe_key
			FROM
				%[1]v
			WHERE recipe_id=$1 AND user_id=$2
		`, usersRecipesTable)

	if err := r.db.Get(&key, setUserPublicKeyQuery, recipeId, userId); err != nil || key == nil {
		if err != nil {
			logRepoError(err)
		}
		return "", failure.NoKey
	}

	return *key, nil
}

func (r *RecipeSharingPostgres) SetOwnerPrivateKeyLinkForUser(recipeId string, userId string, recipeKey *string) error {

	query := fmt.Sprintf(`
			UPDATE %s
			SET recipe_key=$1
			WHERE recipe_id=$2 AND user_id=$3
		`, usersRecipesTable)

	res, err := r.db.Exec(query, recipeKey, recipeId, userId)
	if err != nil {
		logRepoError(err)
		return failure.RecipeNotInRecipeBook
	}

	if rowsCount, err := res.RowsAffected(); err == nil && rowsCount == 0 {
		return failure.RecipeNotInRecipeBook
	}

	return nil
}
