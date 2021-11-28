package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/models"
	"time"
)

type ShoppingList struct {
	db *sqlx.DB
}

func NewShoppingListPostgres(db *sqlx.DB) *ShoppingList {
	return &ShoppingList{db: db}
}

func (r *ShoppingList) GetShoppingList(userId int) (models.ShoppingList, error) {
	var shoppingList models.ShoppingList
	var shoppingListJSON []byte
	query := fmt.Sprintf("SELECT shopping_list FROM %s WHERE user_id=$1", shoppingListTable)
	err := r.db.Get(&shoppingListJSON, query, userId)
	if err != nil {
		return models.ShoppingList{}, err
	}
	err = json.Unmarshal(shoppingListJSON, &shoppingList)
	return shoppingList, err
}

func (r *ShoppingList) SetShoppingList(shoppingList models.ShoppingList, userId int) error {
	shoppingList.Timestamp = time.Now()
	var shoppingListJSON, err = json.Marshal(shoppingList)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("UPDATE %s SET shopping_list=$1 WHERE user_id=$2", shoppingListTable)
	_, err = r.db.Exec(query, shoppingListJSON, userId)
	return err
}