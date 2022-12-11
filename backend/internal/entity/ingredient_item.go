package entity

import "github.com/google/uuid"

type IngredientItem struct {
	Id     uuid.UUID
	Text   string
	Amount *int
	Unit   *string
	Link   *string
	Type   string
}
