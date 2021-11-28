package models

import "time"

type ShoppingList struct {
	Purchases []Purchase    `json:"purchases"`
	Timestamp time.Time `json:"timestamp"`
}

type Purchase struct {
	Id          string `json:"purchase_id"`
	Item        string `json:"name"`
	Multiplier  int    `json:"multiplier"`
	IsPurchased bool   `json:"is_purchased"`
}
