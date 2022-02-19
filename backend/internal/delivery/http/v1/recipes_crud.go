package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/model"
	"net/http"
)

func (h *Handler) getRecipes(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	params := getRecipesRequestParamsByCtx(c)
	params.UserId = userId

	recipes, err := h.services.RecipesCrud.GetRecipesInfoByRequest(*params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, recipes)
}

func (h *Handler) createRecipe(c *gin.Context) {
	recipe, err := getRecipeByCtx(c)
	if err != nil {
		return
	}

	id, err := 	h.services.RecipesCrud.CreateRecipe(recipe)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newIdResponse(c, id, RespRecipeAdded)
}


func (h *Handler) addRecipeToRecipeBook(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	err = h.services.RecipesCrud.AddRecipeToRecipeBook(recipeId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newMessageResponse(c, RespRecipeAdded)
}

func (h *Handler) getRecipe(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	recipe, err := h.services.RecipesCrud.GetRecipeById(recipeId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func (h *Handler) updateRecipe(c *gin.Context) {
	recipe, err := getRecipeByCtx(c)
	if err != nil {
		return
	}

	if err := h.services.RecipesCrud.UpdateRecipe(recipe); err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	newMessageResponse(c, RespRecipeUpdated)
}

func (h *Handler) deleteRecipe(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	if err := h.services.RecipesCrud.DeleteRecipe(recipeId, userId); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newMessageResponse(c, RespRecipeDeleted)
}

func (h *Handler) getRandomPublicRecipe(c *gin.Context) {
	languages, ok := c.GetQueryArray("language")
	if !ok {
		languages = []string{}
	}

	recipe, err := h.services.RecipesCrud.GetRandomPublicRecipe(languages)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrUnableGetRandomRecipe.Error())
		return
	}

	c.JSON(http.StatusOK, recipe)
}