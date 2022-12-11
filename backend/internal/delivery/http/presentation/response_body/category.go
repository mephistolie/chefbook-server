package response_body

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Category struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Cover *string   `json:"cover"`
}

func NewCategory(category entity.Category) Category {
	return Category{
		Id:    category.Id,
		Name:  category.Name,
		Cover: category.Cover,
	}
}

func NewCategories(entities []entity.Category) []Category {
	categories := make([]Category, len(entities))
	for i, category := range entities {
		categories[i] = NewCategory(category)
	}
	return categories
}
