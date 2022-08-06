package dto

import (
	"chefbook-server/internal/entity"
	"time"
)

type ShoppingList struct {
	Purchases []Purchase `json:"purchases"`
	Timestamp time.Time  `json:"timestamp"`
}

func NewShoppingList(shoppingList entity.ShoppingList) ShoppingList {
	purchases := make([]Purchase, len(shoppingList.Purchases))
	for i, purchase := range shoppingList.Purchases {
		purchases[i] = newPurchase(purchase)
	}

	return ShoppingList{
		Purchases: purchases,
		Timestamp: shoppingList.Timestamp,
	}
}

func (l *ShoppingList) Entity() entity.ShoppingList {
	purchases := make([]entity.Purchase, len(l.Purchases))
	for i, purchase := range l.Purchases {
		purchases[i] = purchase.Entity()
	}

	return entity.ShoppingList{
		Purchases: purchases,
		Timestamp: l.Timestamp,
	}
}

type Purchase struct {
	Id          string `json:"purchase_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Multiplier  int    `json:"multiplier,omitempty"`
	IsPurchased bool   `json:"is_purchased"`
}

func newPurchase(purchase entity.Purchase) Purchase {
	return Purchase{
		Id:          purchase.Id,
		Name:        purchase.Name,
		Multiplier:  purchase.Multiplier,
		IsPurchased: purchase.IsPurchased,
	}
}

func (l *Purchase) Entity() entity.Purchase {
	return entity.Purchase{
		Id:          l.Id,
		Name:        l.Name,
		Multiplier:  l.Multiplier,
		IsPurchased: l.IsPurchased,
	}
}
