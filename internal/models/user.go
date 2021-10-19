package models

type User struct {
	Id       int    `json:"-"`
	Email    string `json:"email" binding:"required,email,max=128"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}
