package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type ShoppingList interface {
	GetShoppingList(userId uuid.UUID) (entity.ShoppingList, error)
	SetShoppingList(shoppingList entity.ShoppingList, userId uuid.UUID) error
}
