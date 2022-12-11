package service

import (
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

func (s *ShoppingListService) GetShoppingList(userId string) (entity.ShoppingList, error) {
	return s.repo.GetShoppingList(userId)
}

func (s *ShoppingListService) SetShoppingList(purchases []entity.Purchase, userId string) error {
	shoppingList := entity.ShoppingList{
		Purchases: purchases,
		Timestamp: time.Now().UTC(),
	}

	return s.repo.SetShoppingList(shoppingList, userId)
}

func (s *ShoppingListService) AddToShoppingList(newPurchases []entity.Purchase, userId string) error {
	shoppingList, err := s.repo.GetShoppingList(userId)
	if err != nil {
		return err
	}

	shoppingList.Purchases = append(shoppingList.Purchases, newPurchases...)
	shoppingList.Timestamp = time.Now().UTC()

	return s.repo.SetShoppingList(shoppingList, userId)
}
