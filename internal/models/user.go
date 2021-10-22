package models

import (
	"database/sql"
	"github.com/google/uuid"
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
