package service

import "github.com/mephistolie/chefbook-server/internal/entity"

type ShoppingList interface {
	GetShoppingList(userId uuid.UUID) (entity.ShoppingList, error)
	SetShoppingList(purchases []entity.Purchase, userId uuid.UUID) error
	AddToShoppingList(newPurchases []entity.Purchase, userId uuid.UUID) error
}
