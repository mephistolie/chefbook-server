package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id             int            `json:"-,omitempty" db:"user_id"`
	Email          string         `json:"email" binding:"required,email,max=128"`
	Name           sql.NullString `json:"name,omitempty"`
	Password       string         `json:"password,omitempty" binding:"required,min=8,max=64"`
	IsActivated    bool           `json:"is_activated,omitempty" db:"is_activated"`
	ActivationLink uuid.UUID      `json:"activation_link,omitempty" db:"activation_link"`
	Avatar         sql.NullString `json:"avatar,omitempty"`
	VkId           sql.NullString `json:"vk_id,omitempty" db:"vk_id"`
	Premium        sql.NullTime   `json:"premium,omitempty"`
}

type AuthData struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Session struct {
	RefreshToken string    `json:"refreshToken" bson:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt" bson:"expiresAt"`
}
