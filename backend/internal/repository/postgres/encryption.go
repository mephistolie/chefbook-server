package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type EncryptionPostgres struct {
	db *sqlx.DB
}

func NewEncryptionPostgres(db *sqlx.DB) *EncryptionPostgres {
	return &EncryptionPostgres{
		db: db,
	}
}

func (r *EncryptionPostgres) GetUserKey(userId int) (string, error) {
	var key sql.NullString
	query := fmt.Sprintf("SELECT key FROM %s WHERE user_id=$1", usersTable)
	err := r.db.Get(&key, query, userId)
	return key.String, err
}

func (r *EncryptionPostgres) SetUserKey(userId int, url string) error {
	var key interface{}
	if url != "" { key = url} else { key = nil}
	query := fmt.Sprintf("UPDATE %s SET key=$1 WHERE user_id=$2", usersTable)
	_, err := r.db.Exec(query, key, userId)
	return err
}

func (r *EncryptionPostgres) GetRecipeKey(recipeId int) (string, error) {
	var key sql.NullString
	query := fmt.Sprintf("SELECT key FROM %s WHERE recipe_id=$1", recipesTable)
	err := r.db.Get(&key, query, recipeId)
	return key.String, err
}

func (r *EncryptionPostgres) SetRecipeKey(recipeId int, url string) error {
	var key interface{}
	if url != "" {
		key = url
	} else {
		key = nil
	}
	query := fmt.Sprintf("UPDATE %s SET key=$1 WHERE recipe_id=$2", recipesTable)
	_, err := r.db.Exec(query, key, recipeId)
	return err
}