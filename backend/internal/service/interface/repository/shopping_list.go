package repository

import "github.com/mephistolie/chefbook-server/internal/entity"

type ShoppingList interface {
	GetShoppingList(userId string) (entity.ShoppingList, error)
	SetShoppingList(shoppingList entity.ShoppingList, userId string) error
}
