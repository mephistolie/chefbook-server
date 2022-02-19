package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type ProfilePostgres struct {
	db *sqlx.DB
}

func NewProfilePostgres(db *sqlx.DB) *ProfilePostgres {
	return &ProfilePostgres{db: db}
}

func (r *ProfilePostgres) SetUsername(userId int, username string) error {
	query := fmt.Sprintf("UPDATE %s SET username=$1 WHERE user_id=$2", usersTable)
	_, err := r.db.Exec(query, username, userId)
	return err
}

func (r *ProfilePostgres) IncreaseBroccoins(userId, broccoins int) error {
	query := fmt.Sprintf("UPDATE %s SET broccoins=broccoins+$1 WHERE user_id=$2", usersTable)
	_, err := r.db.Exec(query, broccoins, userId)
	return err
}

func (r *ProfilePostgres) ReduceBroccoins(userId, broccoins int) error {
	query := fmt.Sprintf("UPDATE %s SET broccoins=broccoins-$1 WHERE user_id=$2", usersTable)
	_, err := r.db.Exec(query, broccoins, userId)
	return err
}

func (r *ProfilePostgres) SetAvatar(userId int, url string) error {
	var avatar interface{}
	if url != "" { avatar = url} else { avatar = nil}
	query := fmt.Sprintf("UPDATE %s SET avatar=$1 WHERE user_id=$2", usersTable)
	_, err := r.db.Exec(query, avatar, userId)
	return err
}

func (r *ProfilePostgres) SetPremiumDate(userId int, expiresAt time.Time) error {
	query := fmt.Sprintf("UPDATE %s SET premium=$1 WHERE user_id=$2", usersTable)
	_, err := r.db.Exec(query, expiresAt, userId)
	return err
}

func (r *ProfilePostgres) SetProfileCreationDate(userId int, creationTimestamp time.Time) error {
	query := fmt.Sprintf("UPDATE %s SET registered=$1 WHERE user_id=$2", usersTable)
	_, err := r.db.Exec(query, creationTimestamp, userId)
	return err
}