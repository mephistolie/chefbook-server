package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/repository/postgres/dto"
	"time"
)

type ShoppingList struct {
	db *sqlx.DB
}

func NewShoppingListPostgres(db *sqlx.DB) *ShoppingList {
	return &ShoppingList{db: db}
}

func (r *ShoppingList) GetShoppingList(userId int) (entity.ShoppingList, error) {
	var shoppingList dto.ShoppingList
	var shoppingListBSON []byte

	getShoppingListQuery := fmt.Sprintf(`
			SELECT shopping_list
			FROM %s
			WHERE user_id=$1
		`, shoppingListTable)

	if err := r.db.Get(&shoppingListBSON, getShoppingListQuery, userId); err != nil {
		logRepoError(err)
		return entity.ShoppingList{}, failure.ShoppingListNotFound
	}

	if err := json.Unmarshal(shoppingListBSON, &shoppingList); err != nil {
		logRepoError(err)
		emptyShoppingList := entity.ShoppingList{
			Timestamp: time.Now(),
		}
		_ = r.SetShoppingList(emptyShoppingList, userId)
		return emptyShoppingList, nil
	}

	return shoppingList.Entity(), nil
}

func (r *ShoppingList) SetShoppingList(shoppingList entity.ShoppingList, userId int) error {
	var shoppingListBSON, err = json.Marshal(dto.NewShoppingList(shoppingList))
	if err != nil {
		logRepoError(err)
		return failure.Unknown
	}

	setShoppingListQuery := fmt.Sprintf(`
			UPDATE %s
			SET shopping_list=$1
			WHERE user_id=$2
		`, shoppingListTable)

	if _, err = r.db.Exec(setShoppingListQuery, shoppingListBSON, userId); err != nil {
		logRepoError(err)
		return failure.ShoppingListNotFound
	}

	return nil
}
