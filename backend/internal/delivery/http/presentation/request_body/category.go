package request_body

import "chefbook-server/internal/entity"

type CategoryInput struct {
	Name  string  `json:"name" binding:"required,min=1,max=50"`
	Cover *string `json:"cover"`
}

func (c *CategoryInput) Entity() entity.CategoryInput {
	return entity.CategoryInput{
		Name:  c.Name,
		Cover: c.Cover,
	}
}
