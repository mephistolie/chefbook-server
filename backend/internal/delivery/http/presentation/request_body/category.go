package request_body

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type CategoryInput struct {
	Id    *uuid.UUID `json:"category_id"`
	Name  string     `json:"name" binding:"required,min=1,max=50"`
	Cover *string    `json:"cover"`
}

func (c *CategoryInput) Entity() entity.CategoryInput {
	return entity.CategoryInput{
		Id:    c.Id,
		Name:  c.Name,
		Cover: c.Cover,
	}
}
