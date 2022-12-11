package entity

type Category struct {
	Id     string
	Name   string
	Cover  *string
	UserId string
}

type CategoryInput struct {
	Name  string  `json:"name"`
	Cover *string `json:"cover" binding:"max=20"`
}
