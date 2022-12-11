package entity

import "github.com/google/uuid"

const (
	TypeStep = "step"
)

type CookingItem struct {
	Id       uuid.UUID
	Text     string
	Link     *string
	Time     *int16
	Pictures *[]string
	Type     string
}
