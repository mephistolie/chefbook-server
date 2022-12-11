package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/service/interface/repository"
	"time"
)

type ShoppingListService struct {
	repo repository.ShoppingList
}

func NewShoppingListService(repo repository.ShoppingList) *ShoppingListService {
	return &ShoppingListService{
		repo: repo,
	}
}

func (s *ShoppingListService) GetShoppingList(userId uuid.UUID) (entity.ShoppingList, error) {
	return s.repo.GetShoppingList(userId)
}

func (s *ShoppingListService) SetShoppingList(purchases []entity.Purchase, userId uuid.UUID) error {
	shoppingList := entity.ShoppingList{
		Purchases: purchases,
		Timestamp: time.Now().UTC(),
	}

	return s.repo.SetShoppingList(shoppingList, userId)
}

func (s *ShoppingListService) AddToShoppingList(newPurchases []entity.Purchase, userId uuid.UUID) error {
	shoppingList, err := s.repo.GetShoppingList(userId)
	if err != nil {
		return err
	}

	var purchasesByIds map[uuid.UUID]*entity.Purchase
	var purchasesByName map[string]*entity.Purchase

	for i := range shoppingList.Purchases {
		purchasesByIds[shoppingList.Purchases[i].Id] = &shoppingList.Purchases[i]
		purchasesByName[shoppingList.Purchases[i].Name] = &shoppingList.Purchases[i]
	}

	for i := range newPurchases {
		id := newPurchases[i].Id
		name := newPurchases[i].Name
		if newPurchases[i].Amount > 0 && purchasesByIds[id] != nil {
			(*purchasesByIds[id]).Amount += newPurchases[i].Amount
		} else if purchasesByName[name] != nil {
			(*purchasesByIds[id]).Multiplier += 1
		} else {
			shoppingList.Purchases = append(shoppingList.Purchases, newPurchases[i])
		}
	}
	shoppingList.Timestamp = time.Now().UTC()

	return s.repo.SetShoppingList(shoppingList, userId)
}
