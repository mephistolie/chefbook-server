package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/app/dependencies/service"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/middleware"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/middleware/response"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/request_body"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/response_body"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/response_body/message"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"strconv"
)

const (
	ParamRecipeId = "recipe_id"

	queryAuthorId    = "author_id"
	queryOwned       = "owned"
	querySaved       = "saved"
	querySearch      = "search"
	querySortBy      = "sort_by"
	queryLanguages   = "language"
	queryPage        = "page"
	queryPageSize    = "page_size"
	queryMinTime     = "min_time"
	queryMaxTime     = "max_time"
	queryMinServings = "min_servings"
	queryMaxServings = "max_servings"
	queryMinCalories = "min_calories"
	queryMaxCalories = "max_calories"
)

type RecipeHandler struct {
	middleware middleware.AuthMiddleware
	service    service.Recipe
}

func NewRecipeCrudHandler(middleware middleware.AuthMiddleware, service service.Recipe) *RecipeHandler {
	return &RecipeHandler{
		middleware: middleware,
		service:    service,
	}
}

// GetRecipes Swagger Documentation
// @Summary Get Recipes
// @Security ApiKeyAuth
// @Tags recipes
// @Description Get recipes by query
// @Accept json
// @Produce json
// @Param author_id query int false "Recipes author ID"
// @Param owned query bool false "Get only those recipes that were created by user"
// @Param saved query bool false "Get only those recipes that saved to user recipe book"
// @Param search query string false "Search recipes with specified name"
// @Param sort_by query string false "Sorting. Acceptable values: 'creation_timestamp', 'update_timestamp', 'likes', 'time', 'servings', 'calories'"
// @Param language query []string false "Recipe language codes"
// @Param page query string false "Page of the result"
// @Param page_size query string false "Page size of the result. Maximum is 50"
// @Param min_time query string false "Minimal recipe cooking time"
// @Param max_time query string false "Maximum recipe cooking time"
// @Param min_servings query string false "Minimal recipe servings"
// @Param max_servings query string false "Maximum recipe servings"
// @Param min_calories query string false "Minimal recipe calories"
// @Param max_calories query string false "Maximum recipe calories"
// @Success 200 {object} []response_body.RecipeInfo
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes [get]
func (r *RecipeHandler) GetRecipes(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	query := r.getRecipesQuery(c)
	if err := query.Validate(userId); err != nil {
		response.Failure(c, err)
		return
	}

	recipes, err := r.service.GetRecipes(query.Entity(), userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Success(c, response_body.NewRecipes(recipes))
}

// GetRecipe Swagger Documentation
// @Summary Get Recipe
// @Security ApiKeyAuth
// @Tags recipes
// @Description Get recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} response_body.Recipe
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id} [get]
func (r *RecipeHandler) GetRecipe(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.middleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	recipe, err := r.service.GetRecipe(recipeId, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Success(c, response_body.NewRecipe(recipe))
}

// GetRandomRecipe Swagger Documentation
// @Summary Get Random Recipe
// @Security ApiKeyAuth
// @Tags recipes
// @Description Get random recipe
// @Accept json
// @Produce json
// @Param language query []string false "Recipe language codes"
// @Success 200 {object} response_body.Recipe
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/random [get]
func (r *RecipeHandler) GetRandomRecipe(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	var languages *[]string = nil
	if parsedLanguages, ok := c.GetQueryArray(queryLanguages); ok {
		languages = &parsedLanguages
	}

	recipe, err := r.service.GetRandomRecipe(languages, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Success(c, response_body.NewRecipe(recipe))
}

// AddRecipeToRecipeBook Swagger Documentation
// @Summary Add Recipe to Recipe Book
// @Security ApiKeyAuth
// @Tags recipes
// @Description Add recipe to user recipe book
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/save [post]
func (r *RecipeHandler) AddRecipeToRecipeBook(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.middleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	err = r.service.AddRecipeToRecipeBook(recipeId, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.RecipeAddedToRecipeBook)
}

// RemoveFromRecipeBook Swagger Documentation
// @Summary Remove Recipe from Recipe Book
// @Security ApiKeyAuth
// @Tags recipes
// @Description Remove recipe from user recipe book
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/save [delete]
func (r *RecipeHandler) RemoveFromRecipeBook(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.middleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	err = r.service.RemoveRecipeFromRecipeBook(recipeId, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.RecipeRemovedFromRecipeBook)
}

// SetRecipeCategories Swagger Documentation
// @Summary Set Recipe Categories
// @Security ApiKeyAuth
// @Tags recipes
// @Description Set user categories for recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Param categories body []int true "Recipe categories"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/categories [put]
func (r *RecipeHandler) SetRecipeCategories(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.middleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	var body request_body.RecipeCategoriesInput
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	err = r.service.SetRecipeCategories(recipeId, body.Categories, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.CategoriesUpdated)
}

// MarkRecipeFavourite Swagger Documentation
// @Summary Add Recipe to Favourites
// @Security ApiKeyAuth
// @Tags recipes
// @Description Add recipe to favourites
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/favourite [put]
func (r *RecipeHandler) MarkRecipeFavourite(c *gin.Context) {
	r.setRecipeFavourite(c, true)
}

// UnmarkRecipeFavourite Swagger Documentation
// @Summary Delete Recipe from Favourites
// @Security ApiKeyAuth
// @Tags recipes
// @Description Delete recipe from favourites
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/favourite [delete]
func (r *RecipeHandler) UnmarkRecipeFavourite(c *gin.Context) {
	r.setRecipeFavourite(c, false)
}

func (r *RecipeHandler) setRecipeFavourite(c *gin.Context, favourite bool) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.middleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	err = r.service.SetRecipeFavourite(recipeId, favourite, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.FavouriteStatusUpdated)
}

// LikeRecipe Swagger Documentation
// @Summary Like Recipe
// @Security ApiKeyAuth
// @Tags recipes
// @Description Like recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/likes [put]
func (r *RecipeHandler) LikeRecipe(c *gin.Context) {
	r.setRecipeLiked(c, true)
}

// UnlikeRecipe Swagger Documentation
// @Summary Unlike recipe
// @Security ApiKeyAuth
// @Tags recipes
// @Description Unlike recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/likes [delete]
func (r *RecipeHandler) UnlikeRecipe(c *gin.Context) {
	r.setRecipeLiked(c, false)
}

func (r *RecipeHandler) setRecipeLiked(c *gin.Context, liked bool) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.middleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	err = r.service.SetRecipeLikeStatus(recipeId, liked, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.RecipeLikeSet)
}

func (r *RecipeHandler) getRecipesQuery(c *gin.Context) *request_body.RecipesQuery {
	var params request_body.RecipesQuery

	if authorId, ok := c.GetQuery(queryAuthorId); ok {
		*params.AuthorId = authorId
	}

	if ownedQuery, ok := c.GetQuery(queryOwned); ok {
		params.Owned = ownedQuery == "true"
	}

	if savedQuery, ok := c.GetQuery(querySaved); ok {
		params.Saved = savedQuery == "true"
	}

	if search, ok := c.GetQuery(querySearch); ok {
		*params.Search = search
	}

	if sortBy, ok := c.GetQuery(querySortBy); ok {
		params.SortBy = sortBy
	}

	if query, ok := c.GetQuery(queryPage); ok {
		if page, err := strconv.Atoi(query); err == nil {
			params.Page = page
		}
	}

	if query, ok := c.GetQuery(queryPageSize); ok {
		if pageSize, err := strconv.Atoi(query); err == nil {
			params.PageSize = pageSize
		}
	}

	if languages, ok := c.GetQueryArray(queryLanguages); ok {
		*params.Languages = languages
	}

	if query, ok := c.GetQuery(queryMinTime); ok {
		if minTime, err := strconv.Atoi(query); err == nil {
			*params.MinTime = minTime
		}
	}

	if query, ok := c.GetQuery(queryMaxTime); ok {
		if maxTime, err := strconv.Atoi(query); err == nil {
			*params.MaxTime = maxTime
		}
	}

	if query, ok := c.GetQuery(queryMinServings); ok {
		if minServings, err := strconv.Atoi(query); err == nil {
			*params.MinServings = minServings
		}
	}

	if query, ok := c.GetQuery(queryMaxServings); ok {
		if maxServings, err := strconv.Atoi(query); err == nil {
			*params.MaxServings = maxServings
		}
	}

	if query, ok := c.GetQuery(queryMinCalories); ok {
		if minCalories, err := strconv.Atoi(query); err == nil {
			*params.MinCalories = minCalories
		}
	}

	if query, ok := c.GetQuery(queryMaxCalories); ok {
		if maxCalories, err := strconv.Atoi(query); err == nil {
			*params.MaxCalories = maxCalories
		}
	}

	return &params
}

func getUserAndRecipeIds(c *gin.Context, middleware middleware.AuthMiddleware) (string, string, error) {
	userId, err := middleware.GetUserId(c)
	if err != nil {
		return "", "", err
	}

	recipeId := c.Param(ParamRecipeId)
	if len(recipeId) == 0 {
		return "", "", failure.InvalidRecipe
	}

	return userId, recipeId, nil
}
