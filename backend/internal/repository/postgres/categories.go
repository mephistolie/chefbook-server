package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/model"
)

type CategoriesPostgres struct {
	db *sqlx.DB
}

func NewCategoriesPostgres(db *sqlx.DB) *CategoriesPostgres {
	return &CategoriesPostgres{
		db: db,
	}
}

func (r *CategoriesPostgres) GetUserCategories(userId int) ([]model.Category, error) {
	var categories []model.Category
	query := fmt.Sprintf("SELECT category_id, name, cover, user_id FROM %s WHERE user_id=$1", categoriesTable)
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return []model.Category{}, err
	}
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Cover, &category.UserId)
		if err != nil {
			return []model.Category{}, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *CategoriesPostgres) AddCategory(category model.Category) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, cover, user_id) values ($1, $2, $3) RETURNING category_id", categoriesTable)
	row := r.db.QueryRow(query, category.Name, category.Cover, category.UserId)
	err := row.Scan(&id)
	return id, err
}

func (r *CategoriesPostgres) GetCategoryById(categoryId int) (model.Category, error) {
	var category model.Category
	query := fmt.Sprintf("SELECT category_id, name, cover, user_id FROM %s WHERE category_id=$12", categoriesTable)
	row := r.db.QueryRow(query, categoryId)
	err := row.Scan(&category.Id, &category.Name, &category.Cover, &category.UserId)
	return category, err
}

func (r *CategoriesPostgres) UpdateCategory(category model.Category) error {
	query := fmt.Sprintf("UPDATE %s SET name=$1, cover=$2 WHERE category_id=$3 AND user_id=$4", categoriesTable)
	_, err := r.db.Exec(query, category.Name, category.Cover, category.Id, category.UserId)
	return err
}

func (r *CategoriesPostgres) DeleteCategory(categoryId, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE category_id=$1 AND user_id=$2", categoriesTable)
	_, err := r.db.Exec(query, categoryId, userId)
	return err
}

func (r *CategoriesPostgres) GetRecipeCategories(recipeId, userId int) ([]model.Category, error) {
	var categories []model.Category
	query := fmt.Sprintf("SELECT category_id, name, cover FROM %s WHERE recipe_id=$1 AND user_id=$2", recipesCategoriesTable)
	rows, err := r.db.Query(query, recipeId, userId)
	if err != nil {
		return []model.Category{}, err
	}
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.UserId, &category.Name, &category.Cover)
		if err != nil {
			return []model.Category{}, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}