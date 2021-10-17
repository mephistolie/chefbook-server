package models

type User struct {
	Id       int    `json:"-"`
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"required"`
}
