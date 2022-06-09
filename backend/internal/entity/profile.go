package entity

import (
	"time"
)

type Profile struct {
	Id                int
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
	Id                int
	Username          *string
	CreationTimestamp time.Time
	Avatar            *string
	PremiumEndDate    *time.Time
	Broccoins         int
}
