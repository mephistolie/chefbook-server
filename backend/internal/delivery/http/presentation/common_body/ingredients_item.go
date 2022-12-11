package common_body

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"strings"
)

type IngredientItem struct {
	Id     uuid.UUID `json:"id" binding:"required,min=10"`
	Text   string    `json:"text" binding:"required,min=1"`
	Type   string    `json:"type" binding:"required"`
	Amount *int      `json:"amount,omitempty"`
	Unit   *string   `json:"unit,omitempty"`
	Link   *string   `json:"link,omitempty"`
}

func (i *IngredientItem) Validate() error {
	i.Type = strings.ToLower(i.Type)
	switch i.Type {
	case entity.TypeIngredient, entity.TypeSection:
		if len(i.Text) > 100 {
			return failure.TooLongIngredientItemText
		} else {
			return nil
		}
	case entity.TypeEncryptedData:
		return nil
	default:
		return failure.InvalidIngredientItemType
	}
}

func (i *IngredientItem) IsEncrypted() bool {
	return i.Type == entity.TypeEncryptedData
}

func (i *IngredientItem) Entity() entity.IngredientItem {
	return entity.IngredientItem{
		Id:     i.Id,
		Text:   i.Text,
		Type:   i.Type,
		Amount: i.Amount,
		Unit:   i.Unit,
		Link:   i.Link,
	}
}

func NewIngredientItem(ingredientItem entity.IngredientItem) IngredientItem {
	return IngredientItem{
		Id:     ingredientItem.Id,
		Text:   ingredientItem.Text,
		Type:   ingredientItem.Type,
		Amount: ingredientItem.Amount,
		Unit:   ingredientItem.Unit,
		Link:   ingredientItem.Link,
	}
}
