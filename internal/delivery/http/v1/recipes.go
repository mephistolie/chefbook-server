package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/models"
	"net/http"
)

func (h *Handler) initRecipesRoutes(api *gin.RouterGroup) {
	recipes := api.Group("/recipes", h.userIdentity)
	{
		recipes.POST("/create", h.createRecipe)
		recipes.GET("/:recipe_id", h.readRecipe)
		recipes.PUT("/:recipe_id", h.updateRecipe)
		recipes.DELETE("/:recipe_id", h.deleteRecipe)

		recipes.GET("/public", h.createRecipe)
	}
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
		"id": id,
		"message": models.RespRecipeAdded,
	})
}

func (h *Handler) readRecipe(c *gin.Context) {

}

func (h *Handler) updateRecipe(c *gin.Context) {

}

func (h *Handler) deleteRecipe(c *gin.Context) {

}