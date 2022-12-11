package entity

type IngredientItem struct {
	Id     string
	Text   string
	Amount *int
	Unit   *string
	Link   *string
	Type   string
}
