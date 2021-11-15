package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/models"
)

type ShoppingList struct {
	db *sqlx.DB
}

func NewShoppingListPostgres(db *sqlx.DB) *ShoppingList {
	return &ShoppingList{db: db}
}

func (r *ShoppingList) GetShoppingList(userId int) ([]models.Purchase, error) {
	var shoppingList []models.Purchase
	var shoppingListJSON []byte
	query := fmt.Sprintf("SELECT shopping_list FROM %s WHERE user_id=$1", shoppingListTable)
	err := r.db.Get(&shoppingListJSON, query, userId)
	if err != nil {
		return []models.Purchase{}, err
	}
	err = json.Unmarshal(shoppingListJSON, &shoppingList)
	return shoppingList, err
}

func (r *ShoppingList) SetShoppingList(shoppingList []models.Purchase, userId int) error {
	var shoppingListJSON, err = json.Marshal(shoppingList)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("UPDATE %s SET shopping_list=$1 WHERE user_id=$2", shoppingListTable)
	_, err = r.db.Exec(query, shoppingListJSON, userId)
	return err
}