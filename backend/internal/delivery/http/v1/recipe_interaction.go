package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/model"
	"net/http"
)

func (h *Handler) setRecipeCategories(c *gin.Context) {
	var input model.RecipeCategoriesInput
	var err error
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidInput.Error())
		return
	}
	input.UserId, input.RecipeId, err = getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	err = h.services.RecipeInteraction.SetRecipeCategories(input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrRecipeNotInRecipeBook.Error())
		return
	}

	newMessageResponse(c, RespCategoriesUpdated)
}

func (h *Handler) markRecipeFavourite(c *gin.Context) {
	h.setRecipeFavourite(c, true)
}

func (h *Handler) unmarkRecipeFavourite(c *gin.Context) {
	h.setRecipeFavourite(c, false)
}

func (h *Handler) setRecipeFavourite(c *gin.Context, favourite bool) {
	var err error
	input := model.FavouriteRecipeInput {
		Favourite: favourite,
	}
	input.UserId, input.RecipeId, err = getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	err = h.services.RecipeInteraction.SetRecipeFavourite(input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newMessageResponse(c, RespFavouriteStatusUpdated)
}

func (h *Handler) likeRecipe(c *gin.Context) {
	h.setRecipeLiked(c, true)
}

func (h *Handler) unlikeRecipe(c *gin.Context) {
	h.setRecipeLiked(c, false)
}

func (h *Handler) setRecipeLiked(c *gin.Context, liked bool) {
	var err error
	input := model.RecipeLikeInput {
		Liked: liked,
	}
	input.UserId, input.RecipeId, err = getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	err = h.services.RecipeInteraction.SetRecipeLiked(input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newMessageResponse(c, RespRecipeLikeSet)
}