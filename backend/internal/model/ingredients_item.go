package model

type IngredientsItem struct {
	Text   string `json:"text"`
	Amount int    `json:"amount,omitempty"`
	Unit   string `json:"unit,omitempty"`
	Link   string `json:"link,omitempty"`
	Type   string `json:"type,omitempty"`
}
