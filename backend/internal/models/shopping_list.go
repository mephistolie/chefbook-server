package models

type Purchase struct {
	Item        string  `json:"item"`
	IsPurchased bool `json:"is_selected"`
}
