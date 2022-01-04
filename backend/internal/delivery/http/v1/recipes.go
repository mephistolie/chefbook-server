package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/models"
	"net/http"
	"strconv"
)

func (h *Handler) initRecipesRoutes(api *gin.RouterGroup) {
	recipes := api.Group("/recipes", h.userIdentity)
	{
		recipes.GET("", h.getRecipes)
		recipes.POST("", h.createRecipe)
		recipes.GET("/:recipe_id", h.getRecipe)
		recipes.PUT("/:recipe_id", h.updateRecipe)
		recipes.DELETE("/:recipe_id", h.deleteRecipe)

		recipes.PUT("/:recipe_id/categories", h.setRecipeCategories)
		recipes.PUT("/favourites/:recipe_id", h.markRecipeFavourite)
		recipes.DELETE("/favourites/:recipe_id", h.unmarkRecipeFavourite)
		recipes.PUT("/liked/:recipe_id", h.likeRecipe)
		recipes.DELETE("/liked/:recipe_id", h.unlikeRecipe)
	}
}

func (h *Handler) getRecipes(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	recipes, err := h.services.GetRecipesByUser(userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, recipes)
}

func (h *Handler) createRecipe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var recipe models.Recipe

	if err := c.BindJSON(&recipe); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	recipe.OwnerId = userId
	id, err := 	h.services.Recipes.AddRecipe(recipe)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"message": models.RespRecipeAdded,
	})
}

func (h *Handler) getRecipe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	recipe, err := h.services.GetRecipeById(recipeId, userId)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if recipe.OwnerId != userId && recipe.Visibility == "private" {
		newResponse(c, http.StatusForbidden, models.ErrAccessDenied.Error())
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func (h *Handler) updateRecipe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	var recipe models.Recipe

	if err := c.BindJSON(&recipe); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	recipe.Id = recipeId

	if err := h.services.Recipes.UpdateRecipe(recipe, userId); err != nil {
		newResponse(c, http.StatusForbidden, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespRecipeUpdated,
	})
}

func (h *Handler) deleteRecipe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	if err := h.services.DeleteRecipe(recipeId, userId); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespRecipeDeleted,
	})
}

func (h *Handler) setRecipeCategories(c *gin.Context) {
	var input models.RecipeCategoriesInput
	var err error
	input.UserId, err = getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	input.RecipeId, err = strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.SetRecipeCategories(input)
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrRecipeNotInRecipeBook.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespCategoriesUpdated,
	})
}

func (h *Handler) markRecipeFavourite(c *gin.Context) {
	var err error
	input := models.FavouriteRecipeInput {
		Favourite: true,
	}
	input.UserId, err = getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.RecipeId, err = strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.MarkRecipeFavourite(input)
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrRecipeNotInRecipeBook.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespFavouriteStatusUpdated,
	})
}

func (h *Handler) unmarkRecipeFavourite(c *gin.Context) {
	var err error
	input := models.FavouriteRecipeInput {
		Favourite: false,
	}
	input.UserId, err = getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.RecipeId, err = strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.MarkRecipeFavourite(input)
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrRecipeNotInRecipeBook.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespFavouriteStatusUpdated,
	})
}

func (h *Handler) likeRecipe(c *gin.Context) {
	var err error
	input := models.RecipeLikeInput {
		Liked: true,
	}
	input.UserId, err = getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.RecipeId, err = strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.SetRecipeLike(input)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespRecipeLikeSet,
	})
}

func (h *Handler) unlikeRecipe(c *gin.Context) {
	var err error
	input := models.RecipeLikeInput {
		Liked: true,
	}
	input.UserId, err = getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.RecipeId, err = strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.SetRecipeLike(input)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespRecipeLikeSet,
	})
}