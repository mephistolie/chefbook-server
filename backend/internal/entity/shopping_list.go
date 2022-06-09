package entity

import "time"

type ShoppingList struct {
	Purchases []Purchase
	Timestamp time.Time
}

type Purchase struct {
	Id          string
	Name        string
	Multiplier  int
	IsPurchased bool
}
