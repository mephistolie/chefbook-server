package models

type Purchase struct {
	Id          string `json:"purchase_id"`
	Item        string `json:"name"`
	IsPurchased bool   `json:"is_purchased"`
}
