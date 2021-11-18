package models

type Category struct {
	Id     int    `json:"id,omitempty" db:"category_id"`
	Name   string `json:"name"`
	Type   int    `json:"type"`
	UserId int    `json:"*,omitempty" db:"user_id"`
}
