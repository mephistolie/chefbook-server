package dto

import "chefbook-server/internal/entity"

type IngredientItem struct {
	Text   string  `json:"text" binding:"required,min=1"`
	Type   string  `json:"type" binding:"required"`
	Amount *int    `json:"amount,omitempty"`
	Unit   *string `json:"unit,omitempty"`
	Link   *string `json:"link,omitempty"`
}

type CookingItem struct {
	Text     string    `json:"text" binding:"required,min=1"`
	Type     string    `json:"type" binding:"required"`
	Link     *string   `json:"link,omitempty"`
	Time     *int16    `json:"time,omitempty"`
	Pictures *[]string `json:"pictures,omitempty"`
}

func (i *IngredientItem) Entity() entity.IngredientItem {
	return entity.IngredientItem{
		Text:   i.Text,
		Type:   i.Type,
		Amount: i.Amount,
		Unit:   i.Unit,
		Link:   i.Link,
	}
}

func NewIngredientItem(ingredientItem entity.IngredientItem) IngredientItem {
	return IngredientItem{
		Text:   ingredientItem.Text,
		Type:   ingredientItem.Type,
		Amount: ingredientItem.Amount,
		Unit:   ingredientItem.Unit,
		Link:   ingredientItem.Link,
	}
}

func NewIngredientsEntity(ingredients []IngredientItem) []entity.IngredientItem {
	entities := make([]entity.IngredientItem, len(ingredients))
	for i, ingredient := range ingredients {
		entities[i] = ingredient.Entity()
	}
	return entities
}

func NewIngredients(entities []entity.IngredientItem) []IngredientItem {
	ingredients := make([]IngredientItem, len(entities))
	for i, ingredient := range entities {
		ingredients[i] = NewIngredientItem(ingredient)
	}
	return ingredients
}

func (i *CookingItem) Entity() entity.CookingItem {
	return entity.CookingItem{
		Text:     i.Text,
		Type:     i.Type,
		Time:     i.Time,
		Pictures: i.Pictures,
		Link:     i.Link,
	}
}

func NewCookingItem(cookingItem entity.CookingItem) CookingItem {
	return CookingItem{
		Text:     cookingItem.Text,
		Type:     cookingItem.Type,
		Time:     cookingItem.Time,
		Pictures: cookingItem.Pictures,
		Link:     cookingItem.Link,
	}
}

func NewCookingEntity(cooking []CookingItem) []entity.CookingItem {
	entities := make([]entity.CookingItem, len(cooking))
	for i, cookingItem := range cooking {
		entities[i] = cookingItem.Entity()
	}
	return entities
}

func NewCooking(entities []entity.CookingItem) []CookingItem {
	cooking := make([]CookingItem, len(entities))
	for i, cookingItem := range entities {
		cooking[i] = NewCookingItem(cookingItem)
	}
	return cooking
}
