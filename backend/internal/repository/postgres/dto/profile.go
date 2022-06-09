package dto

import (
	"github.com/mephistolie/chefbook-server/internal/entity"
	"time"
)

type ProfileInfo struct {
	Id                int        `db:"user_id"`
	Email             string     `db:"email"`
	Username          *string    `db:"username,omitempty"`
	CreationTimestamp time.Time  `db:"registered"`
	Password          string     `db:"password"`
	IsActivated       bool       `db:"is_activated"`
	Avatar            *string    `db:"avatar"`
	PremiumEndDate    *time.Time `db:"premium"`
	Broccoins         int        `db:"broccoins"`
	IsBlocked         bool       `db:"is_blocked"`
	Key               *string    `db:"key"`
}

func (p *ProfileInfo) Entity() entity.Profile {
	return entity.Profile{
		Id:                p.Id,
		Email:             p.Email,
		Username:          p.Username,
		CreationTimestamp: p.CreationTimestamp,
		Password:          p.Password,
		IsActivated:       p.IsActivated,
		Avatar:            p.Avatar,
		PremiumEndDate:    p.PremiumEndDate,
		Broccoins:         p.Broccoins,
		IsBlocked:         p.IsBlocked,
	}
}
