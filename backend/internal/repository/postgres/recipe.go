package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/repository/postgres/dto"
)

type RecipePostgres struct {
	db *sqlx.DB
}

func NewRecipePostgres(db *sqlx.DB) *RecipePostgres {
	return &RecipePostgres{
		db: db,
	}
}

func (r *RecipePostgres) GetRecipes(params entity.RecipesQuery, userId uuid.UUID) ([]entity.RecipeInfo, error) {
	var recipes []entity.RecipeInfo
	var rows *sql.Rows
	var err error

	parametrizedQuery := r.getRecipesByParamsQuery(params, userId)

	if params.Search != nil {
		rows, err = r.db.Query(parametrizedQuery, *params.Search)
	} else {
		rows, err = r.db.Query(parametrizedQuery)
	}
	if err != nil {
		logRepoError(err)
		return []entity.RecipeInfo{}, nil
	}

	for rows.Next() {
		var recipe entity.RecipeInfo
		err := rows.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Language, &recipe.Likes, &recipe.Servings,
			&recipe.Time, &recipe.Calories, &recipe.Preview, &recipe.Visibility, &recipe.IsEncrypted, &recipe.CreationTimestamp,
			&recipe.UpdateTimestamp, &recipe.IsFavourite, &recipe.IsLiked, &recipe.OwnerName, &recipe.IsSaved)
		if err != nil {
			logRepoError(err)
			continue
		}
		recipe.IsOwned = recipe.OwnerId == userId
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (r *RecipePostgres) GetRecipe(recipeId uuid.UUID) (entity.Recipe, error) {
	var recipe entity.Recipe
	var bsonIngredients []byte
	var bsonCooking []byte

	getRecipeQuery := fmt.Sprintf(`
			SELECT
				recipe_id, name, owner_id, language, description, likes, servings, time, calories, protein, fats,
				carbohydrates, ingredients, cooking, preview, visibility, encrypted, creation_timestamp, update_timestamp
			FROM
				%s
			WHERE
				recipe_id=$1
		`, recipesTable)

	row := r.db.QueryRow(getRecipeQuery, recipeId)
	if err := row.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Language, &recipe.Description, &recipe.Likes,
		&recipe.Servings, &recipe.Time, &recipe.Calories, &recipe.Macronutrients.Protein, &recipe.Macronutrients.Fats, &recipe.Macronutrients.Carbohydrates,
		&bsonIngredients, &bsonCooking, &recipe.Preview, &recipe.Visibility, &recipe.IsEncrypted, &recipe.CreationTimestamp, &recipe.UpdateTimestamp); err != nil {
		logRepoError(err)
		return entity.Recipe{}, failure.RecipeNotFound
	}

	var ingredients []dto.IngredientItem
	var cooking []dto.CookingItem
	if err := json.Unmarshal(bsonIngredients, &ingredients); err != nil {
		logRepoError(err)
		return entity.Recipe{}, failure.InvalidRecipe
	}
	if err := json.Unmarshal(bsonCooking, &cooking); err != nil {
		logRepoError(err)
		return entity.Recipe{}, failure.InvalidRecipe
	}
	recipe.Ingredients = dto.NewIngredientsEntity(ingredients)
	recipe.Cooking = dto.NewCookingEntity(cooking)

	return recipe, nil
}

func (r *RecipePostgres) GetRandomRecipe(languages *[]string, userId uuid.UUID) (entity.UserRecipe, error) {
	var recipe entity.UserRecipe
	var bsonIngredients []byte
	var bsonCooking []byte

	getRecipeQuery := fmt.Sprintf(`
			SELECT
				%[1]v.recipe_id, %[1]v.name, %[1]v.owner_id, %[1]v.language, %[1]v.description, %[1]v.likes, %[1]v.servings,
				%[1]v.time, %[1]v.calories, %[1]v.protein, %[1]v.fats, %[1]v.carbohydrates, %[1]v.ingredients, %[1]v.cooking,
				%[1]v.preview, %[1]v.visibility, %[1]v.encrypted, %[1]v.creation_timestamp, %[1]v.update_timestamp, 
				coalesce(%[2]v.favourite, false),
				(
					SELECT EXISTS
					(
						SELECT 1
						FROM %[3]v
						WHERE %[3]v.recipe_id=%[1]v.recipe_id AND user_id=$1
					)
				) AS liked, %[4]v.username,
				(
					SELECT EXISTS
					(
						SELECT 1
						FROM %[2]v
						WHERE %[2]v.recipe_id=%[1]v.recipe_id AND user_id=$1
					)
				) AS saved
			FROM
				%[1]v
			LEFT JOIN
				%[2]v ON %[2]v.user_id=$1 AND %[1]v.recipe_id=%[2]v.recipe_id
			LEFT JOIN
				users ON %[4]v.user_id=%[1]v.owner_id
		`, recipesTable, usersRecipesTable, likesTable, usersTable)
	getRecipeQuery += fmt.Sprintf(" WHERE visibility='%s'", entity.VisibilityPublic)
	getRecipeQuery += r.getLanguagesFilter(languages)
	getRecipeQuery += " ORDER BY RANDOM() LIMIT 1"

	row := r.db.QueryRow(getRecipeQuery, userId)
	if err := row.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Language, &recipe.Description, &recipe.Likes,
		&recipe.Servings, &recipe.Time, &recipe.Calories, &recipe.Macronutrients.Protein, &recipe.Macronutrients.Fats,
		&recipe.Macronutrients.Carbohydrates, &bsonIngredients, &bsonCooking, &recipe.Preview, &recipe.Visibility,
		&recipe.IsEncrypted, &recipe.CreationTimestamp, &recipe.UpdateTimestamp, &recipe.IsFavourite, &recipe.IsLiked,
		&recipe.OwnerName, &recipe.IsSaved); err != nil {
		logRepoError(err)
		return entity.UserRecipe{}, failure.UnableGetRandomRecipe
	}

	recipe.IsOwned = recipe.OwnerId == userId

	var ingredients []dto.IngredientItem
	var cooking []dto.CookingItem
	if err := json.Unmarshal(bsonIngredients, &ingredients); err != nil {
		logRepoError(err)
		return entity.UserRecipe{}, failure.InvalidRecipe
	}
	if err := json.Unmarshal(bsonCooking, &cooking); err != nil {
		logRepoError(err)
		return entity.UserRecipe{}, failure.InvalidRecipe
	}
	recipe.Ingredients = dto.NewIngredientsEntity(ingredients)
	recipe.Cooking = dto.NewCookingEntity(cooking)

	return recipe, nil
}

func (r *RecipePostgres) GetRecipeWithUserFields(recipeId, userId uuid.UUID) (entity.UserRecipe, error) {
	var recipe entity.UserRecipe
	var bsonIngredients []byte
	var bsonCooking []byte

	getRecipeQuery := fmt.Sprintf(`
			SELECT
				%[1]v.recipe_id, %[1]v.name, %[1]v.owner_id, %[1]v.language, %[1]v.description, %[1]v.likes, %[1]v.servings,
				%[1]v.time, %[1]v.calories, %[1]v.protein, %[1]v.fats, %[1]v.carbohydrates, %[1]v.ingredients, %[1]v.cooking,
				%[1]v.preview, %[1]v.visibility, %[1]v.encrypted, %[1]v.creation_timestamp, %[1]v.update_timestamp, 
				coalesce(%[2]v.favourite, false),
				(
					SELECT EXISTS
					(
						SELECT 1
						FROM %[3]v
						WHERE %[3]v.recipe_id=%[1]v.recipe_id AND user_id=$1
					)
				) AS liked, %[4]v.username,
				(
					SELECT EXISTS
					(
						SELECT 1
						FROM %[2]v
						WHERE %[2]v.recipe_id=%[1]v.recipe_id AND user_id=$1
					)
				) AS saved
			FROM
				%[1]v
			LEFT JOIN
				%[2]v ON %[2]v.user_id=$1 AND %[1]v.recipe_id=%[2]v.recipe_id
			LEFT JOIN
				users ON %[4]v.user_id=%[1]v.owner_id
			WHERE %[1]v.recipe_id=$2
		`, recipesTable, usersRecipesTable, likesTable, usersTable)

	row := r.db.QueryRow(getRecipeQuery, userId, recipeId)
	if err := row.Scan(&recipe.Id, &recipe.Name, &recipe.OwnerId, &recipe.Language, &recipe.Description, &recipe.Likes, &recipe.Servings,
		&recipe.Time, &recipe.Calories, &recipe.Macronutrients.Protein, &recipe.Macronutrients.Fats, &recipe.Macronutrients.Carbohydrates,
		&bsonIngredients, &bsonCooking, &recipe.Preview, &recipe.Visibility, &recipe.IsEncrypted, &recipe.CreationTimestamp, &recipe.UpdateTimestamp,
		&recipe.IsFavourite, &recipe.IsLiked, &recipe.OwnerName, &recipe.IsSaved); err != nil {
		logRepoError(err)
		return entity.UserRecipe{}, failure.RecipeNotFound
	}

	recipe.IsOwned = recipe.OwnerId == userId

	var ingredients []dto.IngredientItem
	var cooking []dto.CookingItem
	if err := json.Unmarshal(bsonIngredients, &ingredients); err != nil {
		logRepoError(err)
		return entity.UserRecipe{}, failure.InvalidRecipe
	}
	if err := json.Unmarshal(bsonCooking, &cooking); err != nil {
		logRepoError(err)
		return entity.UserRecipe{}, failure.InvalidRecipe
	}
	recipe.Ingredients = dto.NewIngredientsEntity(ingredients)
	recipe.Cooking = dto.NewCookingEntity(cooking)

	return recipe, nil
}

func (r *RecipePostgres) GetRecipeOwnerId(recipeId uuid.UUID) (uuid.UUID, error) {
	var userId uuid.UUID

	getOwnerQuery := fmt.Sprintf(`
			SELECT owner_id
			FROM %s
			WHERE recipe_id=$1
		`, recipesTable)

	err := r.db.Get(&userId, getOwnerQuery, recipeId)
	if err != nil {
		logRepoError(err)
		return uuid.UUID{}, failure.RecipeNotFound
	}

	return userId, err
}

func (r *RecipePostgres) AddRecipeToRecipeBook(recipeId, userId uuid.UUID) error {
	var exists bool

	checkSavingQuery := fmt.Sprintf(`
			SELECT EXISTS
			(
				SELECT 1
				FROM %s
				WHERE recipe_id=$1 AND user_id=$2
			)
		`, usersRecipesTable)

	err := r.db.QueryRow(checkSavingQuery, recipeId, userId).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		logRepoError(err)
		return failure.Unknown
	}

	if exists {
		return nil
	}

	addRecipeQuery := fmt.Sprintf(`
			INSERT INTO %s (recipe_id, user_id)
			VALUES ($1, $2)
		`, usersRecipesTable)

	if _, err := r.db.Exec(addRecipeQuery, recipeId, userId); err != nil {
		logRepoError(err)
		return failure.UnableAddRecipe
	}

	return nil
}

func (r *RecipePostgres) RemoveRecipeFromRecipeBook(recipeId, userId uuid.UUID) error {

	deleteRecipeQuery := fmt.Sprintf(`
			DELETE FROM %s
			WHERE recipe_id=$1 AND user_id=$2
		`, usersRecipesTable)

	if _, err := r.db.Exec(deleteRecipeQuery, recipeId, userId); err != nil {
		logRepoError(err)
		return failure.RecipeNotFound
	}

	return nil
}

func (r *RecipePostgres) SetRecipeCategories(recipeId uuid.UUID, categoriesIds []uuid.UUID, userId uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		logRepoError(err)
		return failure.Unknown
	}

	checkRecipeInRecipeBookQuery := fmt.Sprintf(`
			SELECT EXISTS
			(
				SELECT 1
				FROM %s
				WHERE recipe_id=$1 AND user_id=$2
			)
		`, usersRecipesTable)

	var inRecipeBook bool
	err = r.db.QueryRow(checkRecipeInRecipeBookQuery, recipeId, userId).Scan(&inRecipeBook)
	if err != nil && err != sql.ErrNoRows {
		logRepoError(err)
		return failure.Unknown
	}
	if !inRecipeBook {
		return failure.RecipeNotInRecipeBook
	}

	clearCategoriesQuery := fmt.Sprintf(`
			DELETE FROM %s
			WHERE recipe_id=$1 AND user_id=$2
		`, recipesCategoriesTable)

	if _, err := tx.Exec(clearCategoriesQuery, recipeId, userId); err != nil {
		logRepoError(err)
		if err := tx.Rollback(); err != nil {
			logRepoError(err)
			return failure.Unknown
		}
		return failure.RecipeNotInRecipeBook
	}

	if len(categoriesIds) > 0 {
		categoriesArrayString := "("
		for _, categoryId := range categoriesIds {
			categoriesArrayString += fmt.Sprintf("'%s', ", categoryId)
		}
		categoriesArrayString = categoriesArrayString[:len(categoriesArrayString)-2] + ")"

		addCategoriesQuery := fmt.Sprintf(`
				INSERT INTO %[1]v
					(recipe_id, category_id, user_id)
					SELECT %[2]v.recipe_id, %[3]v.category_id, %[3]v.user_id
					FROM %[3]v
					LEFT JOIN %[2]v ON %[2]v.recipe_id=$1
				WHERE category_id IN %[4]v AND user_id=$2
			`, recipesCategoriesTable, recipesTable, categoriesTable, categoriesArrayString)

		if _, err := tx.Exec(addCategoriesQuery, recipeId, userId); err != nil {
			logRepoError(err)
			return failure.RecipeNotInRecipeBook
		}
	}

	if err := tx.Commit(); err != nil {
		logRepoError(err)
		return failure.Unknown
	}

	return nil
}

func (r *RecipePostgres) SetRecipeFavourite(recipeId uuid.UUID, isFavourite bool, userId uuid.UUID) error {
	query := fmt.Sprintf(`
			UPDATE %s
			SET favourite=$1
			WHERE recipe_id=$2 AND user_id=$3
		`, usersRecipesTable)

	res, err := r.db.Exec(query, isFavourite, recipeId, userId)
	if err != nil {
		logRepoError(err)
		return failure.RecipeNotInRecipeBook
	}

	if changes, err := res.RowsAffected(); err != nil || changes == 0 {
		return failure.RecipeNotInRecipeBook
	}

	return nil
}

func (r *RecipePostgres) SetRecipeLiked(recipeId uuid.UUID, isLiked bool, userId uuid.UUID) error {
	var exists bool

	checkLikeQuery := fmt.Sprintf(`
			SELECT EXISTS
			(
				SELECT 1
				FROM %s
				WHERE recipe_id=$1 AND user_id=$2
			)
		`, likesTable)

	err := r.db.QueryRow(checkLikeQuery, recipeId, userId).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		logRepoError(err)
		return failure.Unknown
	}

	if (exists && isLiked) || (!exists && !isLiked) {
		return nil
	}

	getRecipeVisibilityQuery := fmt.Sprintf(`
			SELECT visibility, owner_id
			FROM %s
			WHERE recipe_id=$1
		`, recipesTable)

	var visibility string
	var ownerId uuid.UUID
	err = r.db.QueryRow(getRecipeVisibilityQuery, recipeId).Scan(&visibility, &ownerId)
	if err != nil {
		logRepoError(err)
		return failure.Unknown
	}
	if visibility == entity.VisibilityPrivate && ownerId != userId {
		return failure.AccessDenied
	}

	tx, err := r.db.Begin()
	if err != nil {
		logRepoError(err)
		return failure.Unknown
	}

	var likeCountQuery string
	if isLiked {

		likeRecipeQuery := fmt.Sprintf(`
				INSERT INTO %s (recipe_id, user_id)
				VALUES ($1, $2)
			`, likesTable)

		if _, err := tx.Exec(likeRecipeQuery, recipeId, userId); err != nil {
			logRepoError(err)
			if err := tx.Rollback(); err != nil {
				logRepoError(err)
				return failure.Unknown
			}
			return failure.Unknown
		}

		likeCountQuery = fmt.Sprintf(`
				UPDATE %s
				SET likes=likes+1
				WHERE recipe_id=$1
			`, recipesTable)

	} else {

		unlikeRecipeQuery := fmt.Sprintf(`
				DELETE FROM %s
				WHERE recipe_id=$1 AND user_id=$2
			`, likesTable)

		if _, err := tx.Exec(unlikeRecipeQuery, recipeId, userId); err != nil {
			logRepoError(err)
			if err := tx.Rollback(); err != nil {
				logRepoError(err)
				return failure.Unknown
			}
			return failure.Unknown
		}

		likeCountQuery = fmt.Sprintf(`
				UPDATE %s
				SET likes=likes-1
				WHERE recipe_id=$1
			`, recipesTable)

	}

	if _, err := tx.Exec(likeCountQuery, recipeId); err != nil {
		logRepoError(err)
		if err := tx.Rollback(); err != nil {
			logRepoError(err)
			return failure.Unknown
		}
		return failure.Unknown
	}

	if err := tx.Commit(); err != nil {
		logRepoError(err)
		return failure.Unknown
	}

	return nil
}

func (r *RecipePostgres) getRecipesByParamsQuery(params entity.RecipesQuery, userId uuid.UUID) string {

	getRecipesQuery := fmt.Sprintf(`
			SELECT
				%[1]v.recipe_id, %[1]v.name, %[1]v.owner_id, %[1]v.language, %[1]v.likes, %[1]v.servings, %[1]v.time,
				%[1]v.calories, %[1]v.preview, %[1]v.visibility, %[1]v.encrypted, %[1]v.creation_timestamp,
				%[1]v.update_timestamp, coalesce(%[2]v.favourite, false),
				(
					SELECT EXISTS
					(
						SELECT 1
						FROM %[3]v
						WHERE %[3]v.recipe_id=%[1]v.recipe_id AND user_id=%[5]v
					)
				) AS liked, %[4]v.username,
				(
					SELECT EXISTS
					(
						SELECT 1
						FROM %[2]v
						WHERE %[2]v.recipe_id=%[1]v.recipe_id AND user_id=%[5]v
					)
				) AS saved
			FROM
				%[1]v
			LEFT JOIN
				%[2]v ON %[2]v.recipe_id=%[1]v.recipe_id
			LEFT JOIN
				%[4]v ON %[4]v.user_id=%[1]v.owner_id
		`, recipesTable, usersRecipesTable, likesTable, usersTable, userId)

	whereStatement := r.getWhereStatement(params, userId)
	pagingStatement := r.getPagingStatement(params)

	return getRecipesQuery + whereStatement + pagingStatement
}

func (r *RecipePostgres) getWhereStatement(params entity.RecipesQuery, userId uuid.UUID) string {
	whereStatement := " WHERE"

	if params.Saved {
		whereStatement += fmt.Sprintf(" %[1]v.user_id=%[2]v AND (%[3]v.owner_id=%[4]v OR %[3]v.visibility<>'%[5]v')",
			usersRecipesTable, userId, recipesTable, userId, entity.VisibilityPrivate)
	} else {
		whereStatement += fmt.Sprintf(" %[1]v.visibility='%[2]v' AND %[1]v.encrypted=false", recipesTable, entity.VisibilityPublic)
	}

	if params.AuthorId != nil {
		whereStatement += fmt.Sprintf(" AND %s.owner_id=%s", recipesTable, *params.AuthorId)
	}

	whereStatement += r.getLanguagesFilter(params.Languages)

	if params.Search != nil {
		whereStatement += fmt.Sprintf(" AND %s.name LIKE ", recipesTable) + "'%' || $1 || '%'"
	}

	whereStatement += r.getRecipesRangeFilter("time", params.MinTime, params.MaxTime)
	whereStatement += r.getRecipesRangeFilter("servings", params.MinServings, params.MaxServings)
	whereStatement += r.getRecipesRangeFilter("calories", params.MinCalories, params.MaxCalories)

	return whereStatement
}

func (r *RecipePostgres) getPagingStatement(params entity.RecipesQuery) string {
	pagingStatement := fmt.Sprintf(" ORDER BY %s", params.SortBy)

	switch params.SortBy {
	case entity.SortingTime, entity.SortingCalories:
		pagingStatement += " ASC"
	default:
		pagingStatement += " DESC"
	}
	pagingStatement += fmt.Sprintf(" LIMIT %d OFFSET %d", params.PageSize, (params.Page-1)*params.PageSize)

	return pagingStatement
}

func (r *RecipePostgres) getRecipesRangeFilter(field string, min, max *int) string {
	filter := ""
	if min != nil {
		filter += fmt.Sprintf(" AND %s.%s>=%d", recipesTable, field, min)
	}
	if max != nil {
		filter += fmt.Sprintf(" AND %s.%s<=%d", recipesTable, field, max)
	}
	return filter
}

func (r *RecipePostgres) getLanguagesFilter(languages *[]string) string {
	filter := ""
	if languages != nil && len(*languages) > 0 {
		filter += " AND language IN ("
		for _, language := range *languages {
			filter += fmt.Sprintf("'%s', ", language)
		}
		filter = filter[0:len(filter)-2] + ")"
	}
	return filter
}
