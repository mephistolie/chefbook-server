package service

import (
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository"
)

type ShoppingListService struct {
	repo repository.ShoppingList
}

func NewShoppingListService(repo repository.ShoppingList) *ShoppingListService {
	return &ShoppingListService{
		repo: repo,
	}
}

func (s *ShoppingListService) GetShoppingList(userId int) ([]models.Purchase, error)  {
	return s.repo.GetShoppingList(userId)
}

func (s *ShoppingListService) SetShoppingList(shoppingList []models.Purchase, userId int) error  {
	return s.repo.SetShoppingList(shoppingList, userId)
}

func (s *ShoppingListService) AddToShoppingList(newPurchases []models.Purchase, userId int) error {
	originalPurchases, err := s.repo.GetShoppingList(userId)
	if err != nil {
		return err
	}
	shoppingList := append(originalPurchases, newPurchases...)
	return s.repo.SetShoppingList(shoppingList, userId)
}