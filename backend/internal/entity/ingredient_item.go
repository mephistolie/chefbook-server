package entity

const (
	TypeIngredient = "ingredient"
)

type IngredientItem struct {
	Text   string
	Amount *int
	Unit   *string
	Link   *string
	Type   string
}
