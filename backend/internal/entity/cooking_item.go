package entity

const (
	TypeStep = "step"
)

type CookingItem struct {
	Id       string
	Text     string
	Link     *string
	Time     *int16
	Pictures *[]string
	Type     string
}
