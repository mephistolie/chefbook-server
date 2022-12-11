package entity

import (
	"github.com/google/uuid"
	"time"
)

const (
	TypeStandard = "standard"
)

type ShoppingList struct {
	Purchases []Purchase
	Timestamp time.Time
}

type Purchase struct {
	Id          uuid.UUID
	Name        string
	Multiplier  int
	IsPurchased bool
	Amount      int
	Unit        *string
	RecipeId    *string
	RecipeName  *string
}
