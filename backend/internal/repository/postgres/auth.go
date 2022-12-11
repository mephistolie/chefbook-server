package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/repository/postgres/dto"
	"strings"
	"time"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(credentials entity.Credentials, activationLink uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID

	tx, err := r.db.Begin()
	if err != nil {
		logRepoError(err)
		return uuid.UUID{}, failure.Unknown
	}

	createUserQuery := fmt.Sprintf(`
			INSERT INTO %s (email, username, password)
			VALUES ($1, $2, $3)
			RETURNING user_id
		`, usersTable)

	username := credentials.Email[0:strings.Index(credentials.Email, "@")]
	row := tx.QueryRow(createUserQuery, credentials.Email, username, credentials.Password)
	if err := row.Scan(&id); err != nil {
		if err := tx.Rollback(); err != nil {
			logRepoError(err)
			return uuid.UUID{}, failure.Unknown
		}
		return uuid.UUID{}, failure.UnableCreateProfile
	}

	createActivationLinkQuery := fmt.Sprintf(`
			INSERT INTO %s (activation_link, user_id)
			VALUES ($1, $2)
			RETURNING user_id
		`, activationLinksTable)

	if _, err := tx.Exec(createActivationLinkQuery, activationLink, id); err != nil {
		if err := tx.Rollback(); err != nil {
			logRepoError(err)
			return uuid.UUID{}, failure.Unknown
		}
		return uuid.UUID{}, failure.UnableCreateProfile
	}

	createRoleQuery := fmt.Sprintf(`
			INSERT INTO %s (name, user_id)
			VALUES ($1, $2)
		`, rolesTable)

	if _, err := tx.Exec(createRoleQuery, "credentials", id); err != nil {
		logRepoError(err)
		if err := tx.Rollback(); err != nil {
			logRepoError(err)
			return uuid.UUID{}, failure.Unknown
		}
		return uuid.UUID{}, failure.UnableCreateProfile
	}

	createShoppingListQuery := fmt.Sprintf(`
			INSERT INTO %s (user_id)
			VALUES ($1)
		`, shoppingListTable)

	if _, err := tx.Exec(createShoppingListQuery, id); err != nil {
		logRepoError(err)
		if err := tx.Rollback(); err != nil {
			logRepoError(err)
			return uuid.UUID{}, failure.Unknown
		}
		return uuid.UUID{}, failure.UnableCreateProfile
	}

	if err = tx.Commit(); err != nil {
		logRepoError(err)
		return uuid.UUID{}, failure.Unknown
	}

	return id, nil
}

func (r *AuthPostgres) GetUserById(userId uuid.UUID) (entity.Profile, error) {
	var user dto.ProfileInfo

	getUserQuery := fmt.Sprintf(`
			SELECT user_id, email, username, registered, password, is_activated,avatar, premium, broccoins, is_blocked
			FROM %s
			WHERE user_id=$1
		`, usersTable)

	if err := r.db.Get(&user, getUserQuery, userId); err != nil {
		logRepoError(err)
		return entity.Profile{}, failure.UserNotFound
	}

	return user.Entity(), nil
}

func (r *AuthPostgres) GetUserByEmail(email string) (entity.Profile, error) {
	var user dto.ProfileInfo

	getUserQuery := fmt.Sprintf(`
			SELECT user_id, email, username, password, is_activated, avatar, premium, broccoins, is_blocked
			FROM %s
			WHERE email=$1
		`, usersTable)

	if err := r.db.Get(&user, getUserQuery, email); err != nil {
		logRepoError(err)
		return entity.Profile{}, failure.UserNotFound
	}

	return user.Entity(), nil
}

func (r *AuthPostgres) GetUserByRefreshToken(refreshToken string) (entity.Profile, error) {
	var userId uuid.UUID
	var session entity.Session

	getUserIdQuery := fmt.Sprintf(`
			SELECT user_id, expires_at
			FROM %s
			WHERE refresh_token=$1
		`, sessionsTable)

	row := r.db.QueryRow(getUserIdQuery, refreshToken)
	if err := row.Scan(&userId, &session.ExpiresAt); err != nil {
		logRepoError(err)
		return entity.Profile{}, failure.Unknown
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = r.DeleteSession(refreshToken)
		return entity.Profile{}, failure.SessionExpired
	}

	return r.GetUserById(userId)
}

func (r *AuthPostgres) GetUserActivationLink(userId uuid.UUID) (uuid.UUID, error) {
	var activationLink uuid.UUID

	getActivationLinkQuery := fmt.Sprintf(`
			SELECT activation_link
			FROM %s
			WHERE user_id=$1
		`, activationLinksTable)

	if err := r.db.Get(&activationLink, getActivationLinkQuery, userId); err != nil {
		logRepoError(err)
		return uuid.UUID{}, failure.ActivationLinkNotFound
	}

	return activationLink, nil
}

func (r *AuthPostgres) ActivateProfile(activationLink uuid.UUID) error {

	activateProfileQuery := fmt.Sprintf(`
			UPDATE %s
			SET is_activated=true
			WHERE user_id=
			(
				SELECT user_id
				FROM %s
				WHERE activation_link=$1
			)
		`, usersTable, activationLinksTable)

	if _, err := r.db.Exec(activateProfileQuery, activationLink); err != nil {
		logRepoError(err)
		return failure.InvalidActivationLink
	}

	return nil
}

func (r *AuthPostgres) ChangePassword(userId uuid.UUID, password string) error {
	id := 0

	changePasswordQuery := fmt.Sprintf(`
			UPDATE %s
			SET password=$1
			WHERE user_id=$2 RETURNING user_id
		`, usersTable)

	row := r.db.QueryRow(changePasswordQuery, password, userId)
	if err := row.Scan(&id); err != nil || id == -1 {
		logRepoError(err)
		return failure.UserNotFound
	}

	return nil
}

func (r *AuthPostgres) CreateSession(session entity.Session) error {

	createSessionQuery := fmt.Sprintf(`
			INSERT INTO %s (user_id, refresh_token, ip, expires_at)
			VALUES ($1, $2, $3, $4)
		`, sessionsTable)

	if _, err := r.db.Exec(createSessionQuery, session.UserId, session.RefreshToken, session.Ip, session.ExpiresAt); err != nil {
		return failure.Unknown
	}
	return nil
}

func (r *AuthPostgres) DeleteOldSessions(userId uuid.UUID, sessionsThreshold int) error {

	deleteOldSessionsQuery := fmt.Sprintf(`
				DELETE FROM %[1]v
				WHERE session_id NOT IN
				(
					SELECT session_id
					FROM %[1]v
					WHERE user_id=$1
					ORDER BY created_at DESC
					LIMIT %[2]v
				)
			`, sessionsTable, sessionsThreshold)

	if _, err := r.db.Exec(deleteOldSessionsQuery, userId); err != nil {
		logRepoError(err)
		return failure.Unknown
	}

	return nil
}

func (r *AuthPostgres) UpdateSession(session entity.Session, oldRefreshToken string) error {

	updateSessionQuery := fmt.Sprintf(`
			UPDATE %s
			SET refresh_token=$1, ip=$2, expires_at=$3
			WHERE refresh_token=$4
		`, sessionsTable)

	if _, err := r.db.Exec(updateSessionQuery, session.RefreshToken, session.Ip, session.ExpiresAt, oldRefreshToken); err != nil {
		logRepoError(err)
		return failure.SessionNotFound
	}

	return nil
}

func (r *AuthPostgres) DeleteSession(refreshToken string) error {
	var id = -1

	deleteSessionQuery := fmt.Sprintf(`
			DELETE FROM %s
			WHERE refresh_token=$1
			RETURNING session_id
		`, sessionsTable)

	row := r.db.QueryRow(deleteSessionQuery, refreshToken)
	if err := row.Scan(&id); err != nil {
		logRepoError(err)
		return failure.UnableDeleteSession
	}

	return nil
}
