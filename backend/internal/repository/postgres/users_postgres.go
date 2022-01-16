package postgres

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/models"
	"time"
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

	createUserQuery := fmt.Sprintf("INSERT INTO %s (email, username, password, activation_link) values " +
		"($1, $2, $3, $4) RETURNING user_id", usersTable)
	row := tx.QueryRow(createUserQuery, user.Email, user.Email, user.Password, activationLink)
	if err := row.Scan(&id); err != nil {
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

func (r *AuthPostgres) GetUserById(id int) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT user_id, email, username, password, is_activated, activation_link, avatar, vk_id, " +
		"premium, broccoins, is_blocked FROM %s WHERE user_id=$1", usersTable)
	err := r.db.Get(&user, query, id)
	return user, err
}

func (r *AuthPostgres) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT user_id, email, username, password, is_activated, activation_link, avatar, vk_id, " +
		"premium, broccoins, is_blocked FROM %s WHERE email=$1", usersTable)
	err := r.db.Get(&user, query, email)
	return user, err
}

func (r *AuthPostgres) GetUserByCredentials(email, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT user_id, email, username, password, is_activated, activation_link, avatar, vk_id, " +
		"premium, broccoins, is_blocked FROM %s WHERE email=$1 AND password=$2", usersTable)
	err := r.db.Get(&user, query, email, password)
	return user, err
}

func (r *AuthPostgres) GetByRefreshToken(refreshToken string) (models.User, error) {
	var userId int
	var session models.Session
	query := fmt.Sprintf("SELECT user_id, expires_at FROM %s WHERE refresh_token=$1", sessionsTable)
	row := r.db.QueryRow(query, refreshToken)
	if err := row.Scan(&userId, &session.ExpiresAt); err != nil || session.ExpiresAt.Before(time.Now()) {
		if err := r.DeleteSession(refreshToken); err != nil {
			return models.User{}, err
		}
		return models.User{}, models.ErrSessionExpired
	}
	return r.GetUserById(userId)
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

func (r *AuthPostgres) CreateSession(session models.Session) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, refresh_token, ip, expires_at) values ($1, $2, $3, $4)", sessionsTable)
	if _, err := r.db.Exec(query, session.UserId, session.RefreshToken, session.Ip, session.ExpiresAt); err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) UpdateSession(session models.Session, oldRefreshToken string) error  {
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

func (r *AuthPostgres) checkRefreshToken(userId int, session models.Session, ip string) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, refresh_token, ip, expires_at) values ($1, $2, $3, $4)", sessionsTable)
	if _, err := r.db.Exec(query, userId, session.RefreshToken, ip, session.ExpiresAt); err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) ChangePassword(authData models.AuthData) error {
	var id = -1
	query := fmt.Sprintf("UPDATE %s SET password=$1 WHERE email=$2 RETURNING user_id", usersTable)
	row := r.db.QueryRow(query, authData.Password, authData.Email)
	if err := row.Scan(&id); err == nil {
		return err
	}
	if id == -1 {
		return errors.New("user not found")
	}
	return nil
}

func (r *AuthPostgres) SetUserName(userId int, username string) error {
	query := fmt.Sprintf("UPDATE %s SET username=$1 WHERE user_id=$2", usersTable)
	_, err := r.db.Exec(query, username, userId)
	return err
}

func (r *AuthPostgres) SetUserAvatar(userId int, url string) error {
	var avatar interface{}
	if url != "" { avatar = url} else { avatar = nil}
	query := fmt.Sprintf("UPDATE %s SET avatar=$1 WHERE user_id=$2", usersTable)
	_, err := r.db.Exec(query, avatar, userId)
	return err
}

func (r *AuthPostgres) GetUserKey(userId int) (string, error) {
	var key string
	query := fmt.Sprintf("SELECT rsa FROM %s WHERE WHERE user_id=$2", usersTable)
	err := r.db.Get(&key, query, userId)
	return key, err
}

func (r *AuthPostgres) SetUserKey(userId int, url string) error {
	var key interface{}
	if url != "" { key = url} else { key = nil}
	query := fmt.Sprintf("UPDATE %s SET rsa=$1 WHERE user_id=$2", usersTable)
	_, err := r.db.Exec(query, key, userId)
	return err
}

func (r *AuthPostgres) SetPremiumDate(userId int, expiresAt time.Time) error {
	query := fmt.Sprintf("UPDATE %s SET premium=$1 WHERE user_id=$2", usersTable)
	_, err := r.db.Exec(query, expiresAt, userId)
	return err
}