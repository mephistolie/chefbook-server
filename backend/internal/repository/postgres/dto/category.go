package dto

import (
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Category struct {
	Id     string  `db:"category_id"`
	Name   string  `db:"name"`
	Cover  *string `db:"cover"`
	UserId string  `db:"user_id"`
}

func (c *Category) Entity() entity.Category {
	return entity.Category{
		Id:     c.Id,
		Name:   c.Name,
		Cover:  c.Cover,
		UserId: c.UserId,
	}
}
