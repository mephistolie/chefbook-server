package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id             int            `json:"user_id,omitempty" db:"user_id"`
	Email          string         `json:"email" binding:"required,email,max=128"`
	Username       sql.NullString `json:"username,omitempty"`
	Password       string         `json:"password,omitempty" binding:"required,min=8,max=64"`
	IsActivated    bool           `json:"is_activated,omitempty" db:"is_activated"`
	ActivationLink uuid.UUID      `json:"activation_link,omitempty" db:"activation_link"`
	Avatar         sql.NullString `json:"avatar,omitempty"`
	VkId           sql.NullString `json:"vk_id,omitempty" db:"vk_id"`
	Premium        sql.NullTime   `json:"premium,omitempty"`
	IsBlocked      bool           `json:"is_blocked,omitempty" db:"is_blocked"`
}

type UserInfo struct {
	Id       int       `json:"user_id,omitempty" db:"user_id"`
	Username string    `json:"username,omitempty"`
	Avatar   string    `json:"avatar,omitempty"`
	Premium  time.Time `json:"premium,omitempty"`
}

type UserDetailedInfo struct {
	Id        int            `json:"user_id,omitempty" db:"user_id"`
	Email     string         `json:"email" binding:"required,email,max=128"`
	Username  string `json:"username,omitempty"`
	Avatar    string         `json:"avatar,omitempty"`
	Premium   time.Time      `json:"premium,omitempty"`
	IsBlocked bool           `json:"is_blocked,omitempty" db:"is_blocked"`
}

type AuthData struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshInput struct {
	Token string `json:"refresh_token" binding:"required"`
}

type Session struct {
	UserId       int       `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" bson:"refreshToken" db:"refresh_token"`
	Ip           string    `json:"ip" bson:"ip" db:"ip"`
	ExpiresAt    time.Time `json:"expires_at" bson:"expiresAt" db:"expires_at"`
}
