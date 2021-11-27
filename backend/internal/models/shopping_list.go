package models

type Purchase struct {
	Id          string `json:"purchase_id"`
	Item        string `json:"name"`
	Multiplier  int    `json:"multiplier"`
	IsPurchased bool   `json:"is_purchased"`
}
