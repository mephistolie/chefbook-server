package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/model"
	"time"
)

type RecipesCrudPostgres struct {
	db *sqlx.DB
}

func NewRecipesPostgres(db *sqlx.DB) *RecipesCrudPostgres {
	return &RecipesCrudPostgres{
		db: db,
	}
}

func (r *RecipesCrudPostgres) GetRecipesInfoByRequest(params model.RecipesRequestParams) ([]model.RecipeInfo, error) {
	var recipes []model.RecipeInfo
	query := getRecipesQuery(params)
	var preview sql.NullString
	var rows *sql.Rows
	var err error
	if params.Search != "" {
		rows, err = r.db.Query(query, params.Search)
	} else {
		rows, err = r.db.Query(query)
	}
	if err != nil {
		return []model.RecipeInfo{}, err
	}
	for rows.Next() {
		var recipe model.RecipeInfo
		err := rows.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Language, &recipe.Likes, &recipe.Servings,
			&recipe.Time, &recipe.Calories, &preview, &recipe.Visibility, &recipe.Encrypted, &recipe.CreationTimestamp,
			&recipe.UpdateTimestamp, &recipe.Favourite, &recipe.Liked, &recipe.OwnerName, &recipe.UserTimestamp)
		if err != nil {
			return []model.RecipeInfo{}, err
		}
		recipe.Preview = preview.String
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r *RecipesCrudPostgres) AddRecipeToRecipeBook(recipeId, userId int) error {
	query := fmt.Sprintf("INSERT INTO %s (recipe_id, user_id) VALUES ($1, $2)", usersRecipesTable)
	_, err := r.db.Exec(query, recipeId, userId)
	return err
}

func (r *RecipesCrudPostgres) CreateRecipe(recipe model.Recipe) (int, error) {
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

	createRecipeQuery := fmt.Sprintf("INSERT INTO %s (name, owner_id, language, description, servings, time, calories,"+
		"ingredients, cooking, preview, visibility, encrypted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) "+
		"RETURNING recipe_id",
		recipesTable)
	row := tx.QueryRow(createRecipeQuery, recipe.Name, recipe.OwnerId, recipe.Language, recipe.Description, recipe.Servings, recipe.Time, recipe.Calories, recipe.Ingredients,
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

func (r *RecipesCrudPostgres) GetRecipe(recipeId int) (model.Recipe, error) {
	var recipe model.Recipe
	query := fmt.Sprintf("SELECT recipe_id, name, owner_id, language, description, likes, servings, time, calories, "+
		"ingredients, cooking, preview, visibility, encrypted, creation_timestamp, update_timestamp FROM %s WHERE recipe_id=$1",
		recipesTable)
	var ingredients []byte
	var cooking []byte
	var preview sql.NullString
	row := r.db.QueryRow(query, recipeId)
	if err := row.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Language, &recipe.Description, &recipe.Likes, &recipe.Servings, &recipe.Time, &recipe.Calories, &ingredients,
		&cooking, &preview, &recipe.Visibility, &recipe.Encrypted, &recipe.CreationTimestamp, &recipe.UpdateTimestamp); err != nil {
		return model.Recipe{}, err
	}
	if err := json.Unmarshal(ingredients, &recipe.Ingredients); err != nil {
		return model.Recipe{}, err
	}
	if err := json.Unmarshal(cooking, &recipe.Cooking); err != nil {
		return model.Recipe{}, err
	}
	recipe.Preview = preview.String
	return recipe, nil
}

func (r *RecipesCrudPostgres) GetRecipeWithUserFields(recipeId, userId int) (model.Recipe, error) {
	var recipe model.Recipe
	query := fmt.Sprintf("SELECT %[1]v.recipe_id, %[1]v.name, %[1]v.owner_id, %[1]v.language, %[1]v.description, %[1]v.likes, "+
		"%[1]v.servings, %[1]v.time, %[1]v.calories, %[1]v.ingredients, %[1]v.cooking, %[1]v.preview, %[1]v.visibility, "+
		"%[1]v.encrypted, %[1]v.creation_timestamp, %[1]v.update_timestamp, coalesce(%[2]v.favourite, false), (SELECT "+
		"EXISTS (SELECT 1 FROM %[3]v WHERE %[3]v.recipe_id=%[1]v.recipe_id AND user_id=$1)) as liked, %[4]v.username, " +
		"%[2]v.update_timestamp FROM %[1]v LEFT JOIN %[2]v ON %[2]v.user_id=$1 AND %[1]v.recipe_id=%[2]v.recipe_id " +
		"LEFT JOIN users ON %[4]v.user_id=%[1]v.owner_id WHERE %[1]v.recipe_id=$2", recipesTable, usersRecipesTable,
		likesTable, usersTable)
	var ingredients []byte
	var cooking []byte
	var preview sql.NullString
	row := r.db.QueryRow(query, userId, recipeId)
	if err := row.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Language, &recipe.Description, &recipe.Likes, &recipe.Servings,
		&recipe.Time, &recipe.Calories, &ingredients, &cooking, &preview, &recipe.Visibility, &recipe.Encrypted,
		&recipe.CreationTimestamp, &recipe.UpdateTimestamp, &recipe.Favourite, &recipe.Liked, &recipe.OwnerName,
		&recipe.UserTimestamp); err != nil {
		return model.Recipe{}, err
	}
	if err := json.Unmarshal(ingredients, &recipe.Ingredients); err != nil {
		return model.Recipe{}, err
	}
	if err := json.Unmarshal(cooking, &recipe.Cooking); err != nil {
		return model.Recipe{}, err
	}
	recipe.Preview = preview.String
	return recipe, nil
}

func (r *RecipesCrudPostgres) GetRecipeOwnerId(recipeId int) (int, error) {
	var userId int
	query := fmt.Sprintf("SELECT owner_id FROM %s WHERE recipe_id=$1", recipesTable)
	err := r.db.Get(&userId, query, recipeId)
	return userId, err
}

func (r *RecipesCrudPostgres) UpdateRecipe(recipe model.Recipe) error {

	var preview interface{}
	if recipe.Preview != "" {
		preview = recipe.Preview
	} else {
		preview = nil
	}

	query := fmt.Sprintf("UPDATE %s SET name=$1, language=$2, description=$3, servings=$4, time=$5, calories=$6, ingredients=$7, "+
		"cooking=$8, preview=$9, visibility=$10, encrypted=$11, update_timestamp=$12 WHERE recipe_id=$13",
		recipesTable)
	_, err := r.db.Exec(query, recipe.Name, recipe.Language, recipe.Description, recipe.Servings, recipe.Time, recipe.Calories, recipe.Ingredients,
		recipe.Cooking, preview, recipe.Visibility, recipe.Encrypted, time.Now(), recipe.Id)
	return err
}

func (r *RecipesCrudPostgres) DeleteRecipe(recipeId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE recipe_id=$1", recipesTable)
	_, err := r.db.Exec(query, recipeId)
	return err
}

func (r *RecipesCrudPostgres) DeleteRecipeFromRecipeBook(recipeId, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE recipe_id=$1 AND user_id=$2", usersRecipesTable)
	_, err := r.db.Exec(query, recipeId, userId)
	return err
}

func getRecipesQuery(params model.RecipesRequestParams) string {
	whereStatement := " WHERE"
	if params.Owned {
		whereStatement += fmt.Sprintf(" %s.user_id=%d", usersRecipesTable, params.UserId)
	} else {
		whereStatement += fmt.Sprintf(" %[1]v.visibility='public' AND %[1]v.encrypted=false", recipesTable)
	}

	if !params.Owned && params.AuthorId != 0 {
		whereStatement += fmt.Sprintf(" AND %s.owner_id=%d", recipesTable, params.AuthorId)
	}

	if len(params.Languages) > 0 {
		whereStatement += fmt.Sprintf(" AND %s.language IN (", recipesTable)
		for _, language := range params.Languages {
			whereStatement += fmt.Sprintf("%s, ", language)
		}
		whereStatement = whereStatement[0:len(whereStatement)-2] + ")"
	}

	if params.Search != "" {
		whereStatement += fmt.Sprintf(" AND %s.name LIKE ", recipesTable) + "'%' || $1 || '%'"
	}

	whereStatement += getRecipesRangeFilter("time", params.MinTime, params.MaxTime)
	whereStatement += getRecipesRangeFilter("servings", params.MinServings, params.MaxServings)
	whereStatement += getRecipesRangeFilter("calories", params.MinCalories, params.MaxCalories)

	pagingStatement := fmt.Sprintf(" ORDER BY %s", params.SortBy)
	switch params.SortBy {
	case "time", "calories":
		pagingStatement += " ASC"
	default:
		pagingStatement += " DESC"
	}
	pagingStatement += fmt.Sprintf(" LIMIT %d OFFSET %d", params.PageSize, (params.Page-1) * params.PageSize)

	return fmt.Sprintf("SELECT %[1]v.recipe_id, %[1]v.name, %[1]v.owner_id, %[1]v.language, %[1]v.likes, " +
		"%[1]v.servings, %[1]v.time, %[1]v.calories, %[1]v.preview, %[1]v.visibility, %[1]v.encrypted, " +
		"%[1]v.creation_timestamp, %[1]v.update_timestamp, coalesce(%[2]v.favourite, false), (SELECT EXISTS " +
		"(SELECT 1 FROM %[3]v WHERE %[3]v.recipe_id=%[1]v.recipe_id AND user_id=%[5]v)) as liked, %[4]v.username, " +
		"%[2]v.update_timestamp FROM %[1]v LEFT JOIN %[2]v ON %[2]v.recipe_id=%[1]v.recipe_id LEFT JOIN %[4]v ON " +
		"%[4]v.user_id=%[1]v.owner_id" + whereStatement + pagingStatement, recipesTable, usersRecipesTable, likesTable,
		usersTable, params.UserId)
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