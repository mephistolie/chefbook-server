package service

import "github.com/mephistolie/chefbook-server/internal/entity"

type ShoppingList interface {
	GetShoppingList(userId string) (entity.ShoppingList, error)
	SetShoppingList(purchases []entity.Purchase, userId string) error
	AddToShoppingList(newPurchases []entity.Purchase, userId string) error
}
