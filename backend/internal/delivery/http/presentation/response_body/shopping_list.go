package response_body

import (
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/common_body"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"time"
)

type ShoppingList struct {
	Purchases []common_body.Purchase `json:"purchases"`
	Timestamp string                 `json:"timestamp"`
}

func NewShoppingList(shoppingList entity.ShoppingList) ShoppingList {
	purchases := make([]common_body.Purchase, len(shoppingList.Purchases))
	for i, purchase := range shoppingList.Purchases {
		purchases[i] = common_body.NewPurchase(purchase)
	}

	return ShoppingList{
		Purchases: purchases,
		Timestamp: shoppingList.Timestamp.UTC().Format(time.RFC3339),
	}
}
