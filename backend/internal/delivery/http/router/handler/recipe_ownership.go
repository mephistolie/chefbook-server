package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/app/dependencies/service"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/middleware"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/middleware/response"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/request_body"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/response_body/message"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
)

type OwnedRecipeHandler struct {
	middleware middleware.AuthMiddleware
	service    service.RecipeOwnership
}

func NewOwnedRecipeHandler(middleware middleware.AuthMiddleware, service service.RecipeOwnership) *OwnedRecipeHandler {
	return &OwnedRecipeHandler{
		middleware: middleware,
		service:    service,
	}
}

// CreateRecipe Swagger Documentation
// @Summary Create Recipe
// @Security ApiKeyAuth
// @Tags recipes
// @Description Create new recipe
// @Accept json
// @Produce json
// @Param input body request_body.RecipeInput true "Recipe"
// @Success 200 {object} response_body.Id
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes [post]
func (r *OwnedRecipeHandler) CreateRecipe(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	var body request_body.RecipeInput
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	if err := body.Validate(); err != nil {
		response.Failure(c, err)
		return
	}

	recipeId, err := r.service.CreateRecipe(body.Entity(), userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.NewId(c, recipeId, message.RecipeCreated)
}

// UpdateRecipe Swagger Documentation
// @Summary Update Recipe
// @Security ApiKeyAuth
// @Tags recipes
// @Description Update recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Param input body request_body.RecipeInput true "Recipe"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id} [put]
func (r *OwnedRecipeHandler) UpdateRecipe(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	recipeId := c.Param(ParamRecipeId)
	if len(recipeId) == 0 {
		response.Failure(c, failure.Unknown)
		return
	}

	var body request_body.RecipeInput
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	if err := body.Validate(); err != nil {
		response.Failure(c, err)
		return
	}

	if err := r.service.UpdateRecipe(body.Entity(), recipeId, userId); err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.RecipeUpdated)
}

// DeleteRecipe Swagger Documentation
// @Summary Delete Recipe
// @Security ApiKeyAuth
// @Tags recipes
// @Description Delete recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id} [delete]
func (r *OwnedRecipeHandler) DeleteRecipe(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	recipeId := c.Param(ParamRecipeId)
	if len(recipeId) == 0 {
		response.Failure(c, failure.Unknown)
		return
	}

	if err := r.service.DeleteRecipe(recipeId, userId); err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.RecipeDeleted)
}
