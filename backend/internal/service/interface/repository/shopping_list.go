package repository

import "chefbook-server/internal/entity"

type ShoppingList interface {
	GetShoppingList(userId int) (entity.ShoppingList, error)
	SetShoppingList(shoppingList entity.ShoppingList, userId int) error
}
