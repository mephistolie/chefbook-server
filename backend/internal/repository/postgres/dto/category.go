package dto

import "chefbook-server/internal/entity"

type Category struct {
	Id     int     `db:"category_id"`
	Name   string  `db:"name"`
	Cover  *string `db:"cover"`
	UserId int     `db:"user_id"`
}

func (c *Category) Entity() entity.Category {
	return entity.Category{
		Id:     c.Id,
		Name:   c.Name,
		Cover:  c.Cover,
		UserId: c.UserId,
	}
}
