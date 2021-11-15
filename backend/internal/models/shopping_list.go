package models

type Purchase struct {
	Item        int  `json:"item"`
	IsPurchased bool `json:"is_selected"`
}
