package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/models"
	"time"
)

type RecipesPostgres struct {
	db *sqlx.DB
}

func NewRecipesPostgres(db *sqlx.DB) *RecipesPostgres {
	return &RecipesPostgres{
		db: db,
	}
}

func (r *RecipesPostgres) GetRecipesInfoByRequest(params models.RecipesRequestParams) ([]models.RecipeInfo, error) {
	var recipes []models.RecipeInfo
	query := getRecipesQuery(params)
	var preview sql.NullString
	rows, err := r.db.Query(query)
	if err != nil {
		return []models.RecipeInfo{}, err
	}
	for rows.Next() {
		var recipe models.RecipeInfo
		err := rows.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Likes, &recipe.Servings,
			&recipe.Time, &recipe.Calories, &preview, &recipe.Visibility, &recipe.Encrypted, &recipe.CreationTimestamp,
			&recipe.UpdateTimestamp, &recipe.Favourite, &recipe.Liked, &recipe.OwnerName)
		if err != nil {
			return []models.RecipeInfo{}, err
		}
		recipe.Preview = preview.String
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r *RecipesPostgres) CreateRecipe(recipe models.Recipe) (int, error) {
	var id int
	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}

	var preview interface{}
	if recipe.Preview != "" {
		preview = recipe.Preview
	} else {
		preview = nil
	}

	createRecipeQuery := fmt.Sprintf("INSERT INTO %s (name, owner_id, description, servings, time, calories,"+
		"ingredients, cooking, preview, visibility, encrypted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"+
		"RETURNING recipe_id",
		recipesTable)
	row := tx.QueryRow(createRecipeQuery, recipe.Name, recipe.OwnerId, recipe.Description, recipe.Servings, recipe.Time, recipe.Calories, recipe.Ingredients,
		recipe.Cooking, preview, recipe.Visibility, recipe.Encrypted)
	if err := row.Scan(&id); err != nil {
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, err
	}

	createRecipeLinkQuery := fmt.Sprintf("INSERT INTO %s (user_id, recipe_id) values ($1, $2) RETURNING user_id",
		usersRecipesTable)
	if _, err := tx.Exec(createRecipeLinkQuery, recipe.OwnerId, id); err != nil {
		if err := tx.Rollback(); err != nil {
			return -1, err
		}
		return -1, err
	}
	err = tx.Commit()
	return id, nil
}

func (r *RecipesPostgres) AddRecipeLink(recipeId, userId int) error {
	query := fmt.Sprintf("INSERT INTO %s (recipe_id, user_id) VALUES ($1, $2)", usersRecipesTable)
	_, err := r.db.Exec(query, recipeId, userId)
	return err
}

func (r *RecipesPostgres) SetRecipeCategories(categoriesIds []int, recipeId, userId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	clearCategoriesQuery := fmt.Sprintf("DELETE FROM %s WHERE recipe_id=$1 AND user_id=$2", recipesCategoriesTable)
	if _, err := tx.Exec(clearCategoriesQuery, recipeId, userId); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if len(categoriesIds) > 0 {
		categoriesArrayString := "("
		for _, categoryId := range categoriesIds {
			categoriesArrayString += fmt.Sprintf("%d, ", categoryId)
		}
		categoriesArrayString = categoriesArrayString[:len(categoriesArrayString)-2] + ")"
		addCategoriesQuery := fmt.Sprintf("INSERT INTO %[1]v (recipe_id, category_id, user_id) "+
			"SELECT %[2]v.recipe_id, %[3]v.category_id, %[3]v.user_id FROM %[3]v LEFT JOIN %[2]v ON %[2]v.recipe_id=$1 "+
			"WHERE category_id IN %[4]v AND user_id=$2",
			recipesCategoriesTable, recipesTable, categoriesTable, categoriesArrayString)
		if _, err := tx.Exec(addCategoriesQuery, recipeId, userId); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *RecipesPostgres) GetRecipe(recipeId int) (models.Recipe, error) {
	var recipe models.Recipe
	query := fmt.Sprintf("SELECT recipe_id, name, owner_id, description, likes, servings, time, calories, "+
		"ingredients, cooking, preview, visibility, encrypted, creation_timestamp, update_timestamp FROM %s WHERE recipe_id=$1",
		recipesTable)
	var ingredients []byte
	var cooking []byte
	var preview sql.NullString
	row := r.db.QueryRow(query, recipeId)
	if err := row.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Description, &recipe.Likes, &recipe.Servings, &recipe.Time, &recipe.Calories, &ingredients,
		&cooking, &preview, &recipe.Visibility, &recipe.Encrypted, &recipe.CreationTimestamp, &recipe.UpdateTimestamp); err != nil {
		return models.Recipe{}, err
	}
	if err := json.Unmarshal(ingredients, &recipe.Ingredients); err != nil {
		return models.Recipe{}, err
	}
	if err := json.Unmarshal(cooking, &recipe.Cooking); err != nil {
		return models.Recipe{}, err
	}
	recipe.Preview = preview.String
	return recipe, nil
}

func (r *RecipesPostgres) GetRecipeWithUserFields(recipeId, userId int) (models.Recipe, error) {
	var recipe models.Recipe
	query := fmt.Sprintf("SELECT %[1]v.recipe_id, %[1]v.name, %[1]v.owner_id, %[1]v.description, %[1]v.likes, "+
		"%[1]v.servings, %[1]v.time, %[1]v.calories, %[1]v.ingredients, %[1]v.cooking, %[1]v.preview, %[1]v.visibility, "+
		"%[1]v.encrypted, %[1]v.creation_timestamp, %[1]v.update_timestamp, coalesce(%[2]v.favourite, false), (SELECT "+
		"EXISTS (SELECT 1 FROM %[3]v WHERE %[3]v.recipe_id=%[1]v.recipe_id AND user_id=$1)) as liked, %[4]v.username "+
		"FROM %[1]v LEFT JOIN %[2]v ON %[2]v.user_id=$1 AND %[1]v.recipe_id=%[2]v.recipe_id LEFT JOIN users ON "+
		"%[4]v.user_id=%[1]v.owner_id WHERE %[1]v.recipe_id=$2",
		recipesTable, usersRecipesTable, likesTable, usersTable)
	var ingredients []byte
	var cooking []byte
	var preview sql.NullString
	row := r.db.QueryRow(query, userId, recipeId)
	if err := row.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Description, &recipe.Likes, &recipe.Servings, &recipe.Time, &recipe.Calories, &ingredients,
		&cooking, &preview, &recipe.Visibility, &recipe.Encrypted, &recipe.CreationTimestamp, &recipe.UpdateTimestamp, &recipe.Favourite, &recipe.Liked, &recipe.OwnerName); err != nil {
		return models.Recipe{}, err
	}
	if err := json.Unmarshal(ingredients, &recipe.Ingredients); err != nil {
		return models.Recipe{}, err
	}
	if err := json.Unmarshal(cooking, &recipe.Cooking); err != nil {
		return models.Recipe{}, err
	}
	recipe.Preview = preview.String
	return recipe, nil
}

func (r *RecipesPostgres) GetRecipeOwnerId(recipeId int) (int, error) {
	var userId int
	query := fmt.Sprintf("SELECT owner_id FROM %s WHERE recipe_id=$1", recipesTable)
	err := r.db.Get(&userId, query, recipeId)
	return userId, err
}

func (r *RecipesPostgres) UpdateRecipe(recipe models.Recipe, userId int) error {

	var preview interface{}
	if recipe.Preview != "" {
		preview = recipe.Preview
	} else {
		preview = nil
	}

	query := fmt.Sprintf("UPDATE %s SET name=$1, description=$2, servings=$3, time=$4, calories=$5, ingredients=$6, "+
		"cooking=$7, preview=$8, visibility=$9, encrypted=$10, update_timestamp=$11 WHERE recipe_id=$12 AND owner_id=$13",
		recipesTable)
	_, err := r.db.Exec(query, recipe.Name, recipe.Description, recipe.Servings, recipe.Time, recipe.Calories, recipe.Ingredients,
		recipe.Cooking, preview, recipe.Visibility, recipe.Encrypted, time.Now(), recipe.Id, userId)
	return err
}

func (r *RecipesPostgres) DeleteRecipe(recipeId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE recipe_id=$1", recipesTable)
	_, err := r.db.Exec(query, recipeId)
	return err
}

func (r *RecipesPostgres) DeleteRecipeFromRecipeBook(recipeId, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE recipe_id=$1 AND user_id=$2", usersRecipesTable)
	_, err := r.db.Exec(query, recipeId, userId)
	return err
}

func (r *RecipesPostgres) SetRecipeLike(recipeId, userId int, isLiked bool) error {
	var exists bool
	checkLikeQuery := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE recipe_id=$1 AND user_id=$2)", likesTable)

	err := r.db.QueryRow(checkLikeQuery, recipeId, userId).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if (exists && isLiked) || (!exists && !isLiked) {
		return models.ErrRecipeLikeSetAlready
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var likeCountQuery string
	if isLiked {
		likeRecipeQuery := fmt.Sprintf("INSERT INTO %s (recipe_id, user_id) values ($1, $2)", likesTable)
		if _, err := tx.Exec(likeRecipeQuery, recipeId, userId); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
		likeCountQuery = fmt.Sprintf("UPDATE %s SET likes=likes+1 WHERE recipe_id=$1", recipesTable)
	} else {
		unlikeRecipeQuery := fmt.Sprintf("DELETE FROM %s WHERE recipe_id=$1 AND user_id=$2", likesTable)
		if _, err := tx.Exec(unlikeRecipeQuery, recipeId, userId); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
		likeCountQuery = fmt.Sprintf("UPDATE %s SET likes=likes-1 WHERE recipe_id=$1", recipesTable)
	}

	if _, err := tx.Exec(likeCountQuery, recipeId); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (r *RecipesPostgres) MarkRecipeFavourite(recipeId, userId int, isFavourite bool) error {
	query := fmt.Sprintf("UPDATE %s SET favourite=$1 WHERE recipe_id=$2 AND user_id=$3", usersRecipesTable)
	_, err := r.db.Exec(query, isFavourite, recipeId, userId)
	return err
}

func (r *RecipesPostgres) SetRecipePreview(recipeId int, url string) error {
	var preview interface{}
	if url != "" {
		preview = url
	} else {
		preview = nil
	}
	query := fmt.Sprintf("UPDATE %s SET preview=$1 WHERE recipe_id=$2", recipesTable)
	_, err := r.db.Exec(query, preview, recipeId)
	return err
}

func (r *RecipesPostgres) GetRecipeKey(recipeId int) (string, error) {
	var key sql.NullString
	query := fmt.Sprintf("SELECT key FROM %s WHERE recipe_id=$1", recipesTable)
	err := r.db.Get(&key, query, recipeId)
	return key.String, err
}

func (r *RecipesPostgres) SetRecipeKey(recipeId int, url string) error {
	var key interface{}
	if url != "" {
		key = url
	} else {
		key = nil
	}
	query := fmt.Sprintf("UPDATE %s SET key=$1 WHERE recipe_id=$2", recipesTable)
	_, err := r.db.Exec(query, key, recipeId)
	return err
}

func (r *RecipesPostgres) SetUserPublicKeyForRecipe(recipeId int, userId int, userKey string) error {
	query := fmt.Sprintf("UPDATE %s SET user_key=$1 WHERE recipe_id=$2 AND user_id=$3", usersRecipesTable)
	_, err := r.db.Exec(query, userKey, recipeId, userId)
	return err
}

func (r *RecipesPostgres) SetUserPrivateKeyForRecipe(recipeId int, userId int, recipeKey string) error {
	query := fmt.Sprintf("UPDATE %s SET recipe_key=$1 WHERE recipe_id=$2 AND user_id=$3", usersRecipesTable)
	_, err := r.db.Exec(query, recipeKey, recipeId, userId)
	return err
}

func (r *RecipesPostgres) GetRecipeUserList(recipeId int) ([]models.UserInfo, error) {
	query := fmt.Sprintf("SELECT %[1]v.user_id, %[2]v.username, %[2]v.avatar, %[2]v.premium "+
		"LEFT JOIN %[2]v ON %[2]v.user_id=%[1]v.user_id WHERE %[1]v.user_id=$1", usersRecipesTable, usersTable)
	rows, err := r.db.Query(query, recipeId)
	if err != nil {
		return []models.UserInfo{}, err
	}
	var users []models.UserInfo
	for rows.Next() {
		var user models.UserInfo
		err := rows.Scan(&user.Id, &user.Username, &user.Avatar, &user.Premium)
		if err != nil {
			return []models.UserInfo{}, err
		}
		users = append(users, user)
	}
	return users, nil
}

func getRecipesQuery(params models.RecipesRequestParams) string {
	whereStatement := " WHERE"
	if params.Owned {
		whereStatement += fmt.Sprintf(" %s.user_id=%d", usersRecipesTable, params.UserId)
	} else {
		whereStatement += fmt.Sprintf(" %[1]v.visibility='public' AND %[1]v.encrypted=false", recipesTable)
	}

	if !params.Owned && params.AuthorId != 0 {
		whereStatement += fmt.Sprintf(" AND %s.owner_id=%d", recipesTable, params.AuthorId)
	}

	if params.Search != "" {
		whereStatement += fmt.Sprintf(" AND %s.name LIKE ", recipesTable) + "%" + params.Search + "%"
	}

	whereStatement += getRecipesRangeFilter("time", params.MinTime, params.MaxTime)
	whereStatement += getRecipesRangeFilter("servings", params.MinServings, params.MaxServings)
	whereStatement += getRecipesRangeFilter("calories", params.MinCalories, params.MaxCalories)

	pagingStatement := ""
	if params.LastRecipeId > 0 {
		pagingStatement += fmt.Sprintf(" AND %s.recipe_id < %d", recipesTable, params.LastRecipeId)
	}
	pagingStatement += fmt.Sprintf(" ORDER BY %s DESC", params.SortBy)
	if params.SortBy != "recipe_id" {
		pagingStatement += ", recipe_id DESC"
	}
	pagingStatement += fmt.Sprintf(" LIMIT %d", params.PageSize)

	return fmt.Sprintf("SELECT %[1]v.recipe_id, %[1]v.name, %[1]v.owner_id, %[1]v.likes, "+
		"%[1]v.servings, %[1]v.time, %[1]v.calories, %[1]v.preview, %[1]v.visibility, %[1]v.encrypted,"+
		"%[1]v.creation_timestamp, %[1]v.update_timestamp, coalesce(%[2]v.favourite, false), (SELECT EXISTS "+
		"(SELECT 1 FROM %[3]v WHERE %[3]v.recipe_id=%[1]v.recipe_id AND user_id=%[5]v)) as liked, %[4]v.username "+
		"FROM %[1]v LEFT JOIN %[2]v ON %[2]v.recipe_id=%[1]v.recipe_id LEFT JOIN %[4]v ON %[4]v.user_id=%[1]v.owner_id"+
		whereStatement, recipesTable, usersRecipesTable, likesTable, usersTable, params.UserId)
}

func getRecipesRangeFilter(field string, min, max int) string {
	filter := ""
	if min > 0 {
		filter += fmt.Sprintf(" AND %s.%s>=%d", recipesTable, field, min)
	}
	if max > 0 {
		filter += fmt.Sprintf(" AND %s.%s<=%d", recipesTable, field, max)
	}
	return filter
}
