package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
)

type EncryptionPostgres struct {
	db *sqlx.DB
}

func NewEncryptionPostgres(db *sqlx.DB) *EncryptionPostgres {
	return &EncryptionPostgres{
		db: db,
	}
}

func (r *EncryptionPostgres) GetUserKeyLink(userId string) (*string, error) {
	var link *string

	getKeyLinkQuery := fmt.Sprintf(`
			SELECT key FROM %s
			WHERE user_id=$1
		`, usersTable)

	if err := r.db.Get(&link, getKeyLinkQuery, userId); err != nil || link == nil {
		if err != nil {
			logRepoError(err)
		}
		return nil, failure.NoKey
	}

	return link, nil
}

func (r *EncryptionPostgres) SetUserKeyLink(userId string, url *string) error {

	setKeyQuery := fmt.Sprintf(`
			UPDATE %s
			SET key=$1
			WHERE user_id=$2
		`, usersTable)

	if _, err := r.db.Exec(setKeyQuery, url, userId); err != nil {
		logRepoError(err)
		return failure.UserNotFound
	}

	return nil
}

func (r *EncryptionPostgres) GetRecipeKeyLink(recipeId string) (*string, error) {
	var link *string

	getKeyLinkQuery := fmt.Sprintf(`
			SELECT key
			FROM %s
			WHERE recipe_id=$1
		`, recipesTable)

	if err := r.db.Get(&link, getKeyLinkQuery, recipeId); err != nil || link == nil {
		if err != nil {
			logRepoError(err)
		}
		return nil, failure.NoKey
	}

	return link, nil
}

func (r *EncryptionPostgres) SetRecipeKeyLink(recipeId string, url *string) error {

	setKeyQuery := fmt.Sprintf(`
			UPDATE %s
			SET key=$1
			WHERE recipe_id=$2
		`, recipesTable)

	if _, err := r.db.Exec(setKeyQuery, url, recipeId); err != nil {
		logRepoError(err)
		return failure.UserNotFound
	}

	return nil
}
