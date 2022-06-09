package entity

type Category struct {
	Id     int
	Name   string
	Cover  *string
	UserId int
}

type CategoryInput struct {
	Name  string  `json:"name"`
	Cover *string `json:"cover" binding:"max=20"`
}