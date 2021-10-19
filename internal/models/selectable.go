package models

type Selectable struct {
	Item       string      `json:"item"`
	IsSelected bool   `json:"is_selected"`
}