package common_body

import (
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"strings"
)

type CookingItem struct {
	Id       string    `json:"id" binding:"required,min=10"`
	Text     string    `json:"text"`
	Type     string    `json:"type"`
	Link     *string   `json:"link,omitempty"`
	Time     *int16    `json:"time,omitempty"`
	Pictures *[]string `json:"pictures,omitempty"`
}

func (i *CookingItem) Validate() error {
	i.Type = strings.ToLower(i.Type)
	switch i.Type {
	case entity.TypeStep, entity.TypeSection, entity.TypeEncryptedData:
		return nil
	default:
		return failure.InvalidCookingItemType
	}
}

func (i *CookingItem) IsEncrypted() bool {
	return i.Type == entity.TypeEncryptedData
}

func (i *CookingItem) Entity() entity.CookingItem {
	return entity.CookingItem{
		Id:       i.Id,
		Text:     i.Text,
		Type:     i.Type,
		Time:     i.Time,
		Pictures: i.Pictures,
		Link:     i.Link,
	}
}

func NewCookingItem(cookingItem entity.CookingItem) CookingItem {
	return CookingItem{
		Id:       cookingItem.Id,
		Text:     cookingItem.Text,
		Type:     cookingItem.Type,
		Time:     cookingItem.Time,
		Pictures: cookingItem.Pictures,
		Link:     cookingItem.Link,
	}
}
