package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/model"
	"strings"
	"time"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewUsersPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user model.AuthData, activationLink uuid.UUID) (int, error) {
	var id int
	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}

	username := user.Email[0:strings.Index(user.Email, "@")]
	createUserQuery := fmt.Sprintf("INSERT INTO %s (email, username, password) values " +
		"($1, $2, $3, $4) RETURNING user_id", usersTable)
	row := tx.QueryRow(createUserQuery, user.Email, username, user.Password)
	if err := row.Scan(&id); err != nil {
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, err
	}

	createActivationLinkQuery := fmt.Sprintf("INSERT INTO %s (activation_link, user_id) values " +
		"($1, $2, $3, $4) RETURNING user_id", activationLinksTable)
	if _, err := tx.Exec(createActivationLinkQuery, activationLink, id); err != nil {
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, err
	}

	createRoleQuery := fmt.Sprintf("INSERT INTO %s (name, user_id) values ($1, $2)", rolesTable)
	if _, err := tx.Exec(createRoleQuery, "user", id); err != nil {
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, err
	}

	createShoppingListQuery := fmt.Sprintf("INSERT INTO %s (user_id) values ($1)", shoppingListTable)
	if _, err := tx.Exec(createShoppingListQuery, id); err != nil {
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, err
	}

	err = tx.Commit()
	return id, nil
}

func (r *AuthPostgres) GetUserById(id int) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT user_id, email, username, password, is_activated, avatar, " +
		"premium, broccoins, is_blocked FROM %s WHERE user_id=$1", usersTable)
	err := r.db.Get(&user, query, id)
	return user, err
}

func (r *AuthPostgres) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT user_id, email, username, password, is_activated, avatar, " +
		"premium, broccoins, is_blocked FROM %s WHERE email=$1", usersTable)
	err := r.db.Get(&user, query, email)
	return user, err
}

func (r *AuthPostgres) GetUserByCredentials(email, password string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT user_id, email, username, password, is_activated, avatar, " +
		"premium, broccoins, is_blocked FROM %s WHERE email=$1 AND password=$2", usersTable)
	err := r.db.Get(&user, query, email, password)
	return user, err
}

func (r *AuthPostgres) GetByRefreshToken(refreshToken string) (model.User, error) {
	var userId int
	var session model.Session
	query := fmt.Sprintf("SELECT user_id, expires_at FROM %s WHERE refresh_token=$1", sessionsTable)
	row := r.db.QueryRow(query, refreshToken)
	if err := row.Scan(&userId, &session.ExpiresAt); err != nil || session.ExpiresAt.Before(time.Now()) {
		if err := r.DeleteSession(refreshToken); err != nil {
			return model.User{}, err
		}
		return model.User{}, model.ErrSessionExpired
	}
	return r.GetUserById(userId)
}

func (r *AuthPostgres) GetUserActivationLink(id int) (uuid.UUID, error)  {
	var activationLink uuid.UUID
	query := fmt.Sprintf("SELECT activation_link FROM %s WHERE user_id=$1", activationLinksTable)
	err := r.db.Get(&activationLink, query, id)
	return activationLink, err
}

func (r *AuthPostgres) ActivateUser(activationLink uuid.UUID) error {
	var id = -1
	query := fmt.Sprintf("UPDATE %s SET is_activated=true WHERE user_id=(SELECT user_id from %s WHERE activation_link=$1)", usersTable, activationLinksTable)
	row := r.db.QueryRow(query, activationLink)
	if err := row.Scan(&id); err == nil {
		return err
	}
	if id == -1 {
		return model.ErrInvalidActivationLink
	}
	return nil
}

func (r *AuthPostgres) CreateSession(session model.Session) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, refresh_token, ip, expires_at) values ($1, $2, $3, $4)", sessionsTable)
	if _, err := r.db.Exec(query, session.UserId, session.RefreshToken, session.Ip, session.ExpiresAt); err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) UpdateSession(session model.Session, oldRefreshToken string) error  {
	query := fmt.Sprintf("UPDATE %s SET refresh_token=$1, ip=$2, expires_at=$3 WHERE refresh_token=$4", sessionsTable)
	_, err := r.db.Exec(query, session.RefreshToken, session.Ip, session.ExpiresAt, oldRefreshToken)
	return err
}

func (r *AuthPostgres) DeleteSession(refreshToken string) error {
	var id = -1
	query := fmt.Sprintf("DELETE FROM %s WHERE refresh_token=$1 RETURNING session_id", sessionsTable)
	row := r.db.QueryRow(query, refreshToken)
	return row.Scan(&id)
}

func (r *AuthPostgres) checkRefreshToken(userId int, session model.Session, ip string) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, refresh_token, ip, expires_at) values ($1, $2, $3, $4)", sessionsTable)
	if _, err := r.db.Exec(query, userId, session.RefreshToken, ip, session.ExpiresAt); err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) ChangePassword(authData model.AuthData) error {
	var id = -1
	query := fmt.Sprintf("UPDATE %s SET password=$1 WHERE email=$2 RETURNING user_id", usersTable)
	row := r.db.QueryRow(query, authData.Password, authData.Email)
	if err := row.Scan(&id); err == nil {
		return err
	}
	if id == -1 {
		return model.ErrUserIdNotFound
	}
	return nil
}