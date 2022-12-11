package entity

import (
	"github.com/google/uuid"
	"time"
)

type Profile struct {
	Id                uuid.UUID
	Email             string
	Username          *string
	CreationTimestamp time.Time
	Password          string
	IsActivated       bool
	Avatar            *string
	PremiumEndDate    *time.Time
	Broccoins         int
	IsBlocked         bool
}

type ProfileInfo struct {
	Id                uuid.UUID
	Username          *string
	CreationTimestamp time.Time
	Avatar            *string
	PremiumEndDate    *time.Time
	Broccoins         int
}
