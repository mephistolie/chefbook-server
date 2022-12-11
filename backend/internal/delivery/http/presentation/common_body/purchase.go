package common_body

import "github.com/mephistolie/chefbook-server/internal/entity"

type Purchase struct {
	Id          string  `json:"purchase_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Multiplier  *int    `json:"multiplier,omitempty"`
	IsPurchased bool    `json:"is_purchased"`
	Amount      *int    `json:"amount,omitempty"`
	Unit        *string `json:"unit,omitempty"`
	RecipeId    *string `json:"recipe_id,omitempty"`
	RecipeName  *string `json:"recipe_name,omitempty"`
}

func (l *Purchase) Entity() entity.Purchase {
	multiplier := 1
	if l.Multiplier != nil && *l.Multiplier > 1 {
		multiplier = *l.Multiplier
	}
	amount := 0
	if l.Amount != nil && *l.Amount > 1 {
		amount = *l.Amount
	}
	return entity.Purchase{
		Id:          l.Id,
		Name:        l.Name,
		Multiplier:  multiplier,
		IsPurchased: l.IsPurchased,
		Amount:      amount,
		Unit:        l.Unit,
		RecipeId:    l.RecipeId,
		RecipeName:  l.RecipeName,
	}
}

func NewPurchase(purchase entity.Purchase) Purchase {
	var amount *int = nil
	if purchase.Amount > 0 {
		amount = &purchase.Amount
	}
	return Purchase{
		Id:          purchase.Id,
		Name:        purchase.Name,
		Multiplier:  &purchase.Multiplier,
		IsPurchased: purchase.IsPurchased,
		Amount:      amount,
		Unit:        purchase.Unit,
		RecipeId:    purchase.RecipeId,
		RecipeName:  purchase.RecipeName,
	}
}
