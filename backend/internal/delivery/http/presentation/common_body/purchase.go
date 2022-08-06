package common_body

import "chefbook-server/internal/entity"

type Purchase struct {
	Id          string `json:"purchase_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Multiplier  *int   `json:"multiplier,omitempty"`
	IsPurchased bool   `json:"is_purchased"`
}

func (l *Purchase) Entity() entity.Purchase {
	multiplier := 1
	if l.Multiplier != nil && *l.Multiplier > 1 {
		multiplier = *l.Multiplier
	}
	return entity.Purchase{
		Id:          l.Id,
		Name:        l.Name,
		Multiplier:  multiplier,
		IsPurchased: l.IsPurchased,
	}
}

func NewPurchase(purchase entity.Purchase) Purchase {
	return Purchase{
		Id:          purchase.Id,
		Name:        purchase.Name,
		Multiplier:  &purchase.Multiplier,
		IsPurchased: purchase.IsPurchased,
	}
}
