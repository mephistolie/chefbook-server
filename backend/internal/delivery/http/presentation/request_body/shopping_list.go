package request_body

import (
	"chefbook-server/internal/delivery/http/presentation/common_body"
	"chefbook-server/internal/entity"
)

func NewShoppingListEntity(shoppingList []common_body.Purchase) []entity.Purchase {
	purchases := make([]entity.Purchase, len(shoppingList))
	for i, purchase := range shoppingList {
		purchases[i] = purchase.Entity()
	}

	return purchases
}
