package postgres

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/pkg/logger"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User, activationLink uuid.UUID) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (email, password, activation_link) values ($1, $2, $3) RETURNING user_id", usersTable)
	row := r.db.QueryRow(query, user.Email, user.Password, activationLink)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUserById(id int) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", usersTable)
	err := r.db.Get(&user, query, id)
	return user, err
}

func (r *AuthPostgres) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1", usersTable)
	err := r.db.Get(&user, query, email)
	logger.Error(user.Name.Value())
	return user, err
}

func (r *AuthPostgres) ActivateUser(activationLink uuid.UUID) error {
	var id = -1
	query := fmt.Sprintf("UPDATE %s SET is_activated=true WHERE activation_link=$1 RETURNING user_id", usersTable)
	row := r.db.QueryRow(query, activationLink)
	if err := row.Scan(&id); err == nil {
		return err
	}
	if id == -1 {
		return errors.New("user not found")
	}
	return nil
}