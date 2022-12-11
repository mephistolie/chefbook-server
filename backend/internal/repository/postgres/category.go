package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/repository/postgres/dto"
	"github.com/mephistolie/chefbook-server/pkg/logger"
)

type CategoryPostgres struct {
	db *sqlx.DB
}

func NewCategoryPostgres(db *sqlx.DB) *CategoryPostgres {
	return &CategoryPostgres{
		db: db,
	}
}

func (r *CategoryPostgres) GetUserCategories(userId uuid.UUID) []entity.Category {
	var categories []entity.Category

	getCategoriesQuery := fmt.Sprintf(`
			SELECT category_id, name, cover
			FROM %s
			WHERE user_id=$1
		`, categoriesTable)

	rows, err := r.db.Query(getCategoriesQuery, userId)
	if err != nil {
		return []entity.Category{}
	}

	for rows.Next() {
		category := dto.Category{
			UserId: userId,
		}
		err := rows.Scan(&category.Id, &category.Name, &category.Cover)
		if err != nil {
			logger.Error(err)
			continue
		}
		categories = append(categories, category.Entity())
	}

	return categories
}

func (r *CategoryPostgres) GetRecipeCategories(recipeId, userId uuid.UUID) []entity.Category {
	var categories []entity.Category

	getCategoriesQuery := fmt.Sprintf(`
			SELECT
				%[1]v.category_id, %[2]v.name, %[2]v.cover
			FROM
				%[1]v
			LEFT JOIN
				%[2]v ON %[1]v.user_id= %[2]v.user_id
			WHERE
				%[1]v.recipe_id=$1 AND %[1]v.user_id=$2
		`, recipesCategoriesTable, categoriesTable)

	rows, err := r.db.Query(getCategoriesQuery, recipeId, userId)
	if err != nil {
		logRepoError(err)
		return []entity.Category{}
	}

	for rows.Next() {
		var category dto.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Cover)
		if err != nil {
			continue
		}
		categories = append(categories, category.Entity())
	}

	return categories
}

func (r *CategoryPostgres) CreateCategory(category entity.CategoryInput, userId uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	if category.Id != nil {
		id = *category.Id
	} else {
		id = uuid.New()
	}

	addCategoryQuery := fmt.Sprintf(`
			INSERT INTO %s (id, name, cover, user_id)
			VALUES ($1, $2, $3, $4)
		`, categoriesTable)

	if _, err := r.db.Exec(addCategoryQuery, id, category.Name, category.Cover, userId); err != nil {
		logRepoError(err)
		return uuid.UUID{}, failure.UnableAddCategory
	}

	return id, nil
}

func (r *CategoryPostgres) GetCategory(categoryId uuid.UUID) (entity.Category, error) {
	var category dto.Category

	getCategoryQuery := fmt.Sprintf(`
			SELECT category_id, name, cover, user_id
			FROM %s
			WHERE category_id=$1
		`, categoriesTable)

	row := r.db.QueryRow(getCategoryQuery, categoryId)
	if err := row.Scan(&category.Id, &category.Name, &category.Cover, &category.UserId); err != nil {
		logRepoError(err)
		return entity.Category{}, failure.CategoryNotFound
	}

	return category.Entity(), nil
}

func (r *CategoryPostgres) GetCategoryOwnerId(categoryId uuid.UUID) (uuid.UUID, error) {
	var ownerId uuid.UUID

	getCategoryOwnerIdQuery := fmt.Sprintf(`
			SELECT user_id
			FROM %s
			WHERE category_id=$1
		`, categoriesTable)

	row := r.db.QueryRow(getCategoryOwnerIdQuery, categoryId)
	if err := row.Scan(&ownerId); err != nil {
		logRepoError(err)
		return uuid.UUID{}, failure.CategoryNotFound
	}

	return ownerId, nil
}

func (r *CategoryPostgres) UpdateCategory(categoryId uuid.UUID, category entity.CategoryInput) error {

	updateCategoryQuery := fmt.Sprintf(`
			UPDATE %s
			SET name=$1, cover=$2
			WHERE category_id=$3
		`, categoriesTable)

	if _, err := r.db.Exec(updateCategoryQuery, category.Name, category.Cover, categoryId); err != nil {
		logRepoError(err)
		return failure.CategoryNotFound
	}

	return nil
}

func (r *CategoryPostgres) DeleteCategory(categoryId uuid.UUID) error {

	deleteCategoryQuery := fmt.Sprintf(`
			DELETE FROM %s
			WHERE category_id=$1
		`, categoriesTable)

	if _, err := r.db.Exec(deleteCategoryQuery, categoryId); err != nil {
		logRepoError(err)
		return failure.CategoryNotFound
	}

	return nil
}
