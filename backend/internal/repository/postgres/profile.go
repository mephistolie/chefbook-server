package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"time"
)

type ProfilePostgres struct {
	db *sqlx.DB
}

func NewProfilePostgres(db *sqlx.DB) *ProfilePostgres {
	return &ProfilePostgres{db: db}
}

func (r *ProfilePostgres) SetUsername(userId uuid.UUID, username *string) error {

	SetUsernameQuery := fmt.Sprintf(`
			UPDATE %s
			SET username=$1
			WHERE user_id=$2
		`, usersTable)

	if _, err := r.db.Exec(SetUsernameQuery, username, userId); err != nil {
		logRepoError(err)
		return failure.UserNotFound
	}

	return nil
}

func (r *ProfilePostgres) IncreaseBroccoins(userId uuid.UUID, broccoins int) error {

	increaseBroccoinsQuery := fmt.Sprintf(`
			UPDATE %s
			SET broccoins=broccoins+$1
			WHERE user_id=$2
		`, usersTable)

	if _, err := r.db.Exec(increaseBroccoinsQuery, broccoins, userId); err != nil {
		logRepoError(err)
		return failure.UserNotFound
	}

	return nil
}

func (r *ProfilePostgres) SetAvatarLink(userId uuid.UUID, url *string) error {

	setAvatarQuery := fmt.Sprintf(`
			UPDATE %s
			SET avatar=$1
			WHERE user_id=$2
		`, usersTable)

	if _, err := r.db.Exec(setAvatarQuery, url, userId); err != nil {
		logRepoError(err)
		return failure.UserNotFound
	}

	return nil
}

func (r *ProfilePostgres) SetPremiumDate(userId uuid.UUID, expiresAt time.Time) error {

	setPremiumDateQuery := fmt.Sprintf(`
			UPDATE %s
			SET premium=$1
			WHERE user_id=$2
		`, usersTable)

	if _, err := r.db.Exec(setPremiumDateQuery, expiresAt, userId); err != nil {
		logRepoError(err)
		return failure.UserNotFound
	}
	return nil
}

func (r *ProfilePostgres) SetProfileCreationDate(userId uuid.UUID, creationTimestamp time.Time) error {

	setProfileCreationDate := fmt.Sprintf(`
			UPDATE %s
			SET registered=$1
			WHERE user_id=$2
		`, usersTable)

	if _, err := r.db.Exec(setProfileCreationDate, creationTimestamp, userId); err != nil {
		logRepoError(err)
		return failure.UserNotFound
	}

	return nil
}
