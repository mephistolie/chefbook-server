package entity

const (
	TypeStep = "step"
)

type CookingItem struct {
	Text     string
	Link     *string
	Time     *int16
	Pictures *[]string
	Type     string
}
