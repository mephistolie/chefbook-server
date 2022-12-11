package entity

import (
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
	Id          string
	Type        string
	Name        string
	Multiplier  int
	IsPurchased bool
	Amount      int
	Unit        string
	RecipeId    string
}
