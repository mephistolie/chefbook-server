package models

type Selectable[T interface{}] struct {
	Item       T      `json:"item"`
	IsSelected bool   `json:"is_selected"`
}
