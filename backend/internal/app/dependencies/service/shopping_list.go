package service

import "chefbook-server/internal/entity"

type ShoppingList interface {
	GetShoppingList(userId int) (entity.ShoppingList, error)
	SetShoppingList(purchases []entity.Purchase, userId int) error
	AddToShoppingList(newPurchases []entity.Purchase, userId int) error
}
