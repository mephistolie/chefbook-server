package entity

import "github.com/google/uuid"

type Category struct {
	Id     uuid.UUID
	Name   string
	Cover  *string
	UserId uuid.UUID
}

type CategoryInput struct {
	Id    *uuid.UUID `json:"category_id"`
	Name  string     `json:"name"`
	Cover *string    `json:"cover" binding:"max=20"`
}
