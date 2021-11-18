package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/models"
)

type CategoriesPostgres struct {
	db *sqlx.DB
}

func NewCategoriesPostgres(db *sqlx.DB) *CategoriesPostgres {
	return &CategoriesPostgres{
		db: db,
	}
}

func (r *CategoriesPostgres) GetCategoriesByUser(userId int) ([]models.Category, error) {
	var categories []models.Category
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", categoriesTable)
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return []models.Category{}, err
	}
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Type, &category.UserId)
		if err != nil {
			return []models.Category{}, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *CategoriesPostgres) AddCategory(category models.Category) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, type, user_id) values ($1, $2, $3) RETURNING category_id", categoriesTable)
	row := r.db.QueryRow(query, category.Name, category.Type, category.UserId)
	err := row.Scan(&id)
	return id, err
}

func (r *CategoriesPostgres) GetCategoryById(categoryId, userId int) (models.Category, error) {
	var category models.Category
	query := fmt.Sprintf("SELECT %s WHERE category_id=$1 AND user_id=$2", categoriesTable)
	row := r.db.QueryRow(query, categoryId, userId)
	err := row.Scan(&category.Id, &category.Name, &category.Type, &category.UserId)
	return category, err
}

func (r *CategoriesPostgres) UpdateCategory(category models.Category) error {
	query := fmt.Sprintf("UPDATE %s SET name=$1, type=$2 WHERE category_id=$1 AND user_id=$2", categoriesTable)
	_, err := r.db.Exec(query, category.Id, category.UserId)
	return err
}

func (r *CategoriesPostgres) DeleteCategory(categoryId, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s category_id=$1 AND user_id=$2", categoriesTable)
	_, err := r.db.Exec(query, categoryId, userId)
	return err
}