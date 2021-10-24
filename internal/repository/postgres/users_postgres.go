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

func NewUsersPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.AuthData, activationLink uuid.UUID) (int, error) {
	var id int
	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}

	createUserQuery := fmt.Sprintf("INSERT INTO %s (email, password, activation_link) values ($1, $2, $3) RETURNING user_id", usersTable)
	row := tx.QueryRow(createUserQuery, user.Email, user.Password, activationLink)
	if err := row.Scan(&id); err != nil {
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, err
	}

	createRoleQuery := fmt.Sprintf("INSERT INTO %s (name, user_id) values ($1, $2) RETURNING user_id", rolesTable)
	if _, err := tx.Exec(createRoleQuery, "user", id); err != nil {
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, err
	}

	err = tx.Commit()
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
	return user, err
}

func (r *AuthPostgres) GetUserByCredentials(email, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 AND password=$2", usersTable)
	err := r.db.Get(&user, query, email, password)
	logger.Error(password)
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

func (r *AuthPostgres) CreateSession(userId int, session models.Session, ip string) error {
	logger.Error(userId, session.ExpiresAt, session.RefreshToken, ip)
	query := fmt.Sprintf("INSERT INTO %s (user_id, refresh_token, ip, expires_in) values ($1, $2, $3, $4)", sessionsTable)
	if _, err := r.db.Exec(query, userId, session.RefreshToken, ip, session.ExpiresAt); err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) RefreshTokens()  {

}

func (r *AuthPostgres) ChangePassword(authData models.AuthData) error {
	var id = -1
	query := fmt.Sprintf("UPDATE %s SET password=$1 WHERE email=$2 RETURNING user_id", usersTable)
	row := r.db.QueryRow(query, authData.Password, authData.Password)
	if err := row.Scan(&id); err == nil {
		return err
	}
	if id == -1 {
		return errors.New("user not found")
	}
	return nil
}