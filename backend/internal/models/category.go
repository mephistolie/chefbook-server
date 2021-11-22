package models

type Category struct {
	Id     int    `json:"id,omitempty" db:"category_id"`
	Name   string `json:"name"`
	Cover  string `json:"cover" binding:"max=20"`
	UserId int    `json:"-" db:"user_id"`
}
