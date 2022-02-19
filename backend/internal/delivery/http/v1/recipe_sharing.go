package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/model"
	"net/http"
	"strconv"
)

func (h *Handler) getRecipeUsers(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	userList, err := h.services.RecipeSharing.GetRecipeUserList(recipeId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, userList)
}

func (h *Handler) setRecipePublicKey(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	var keys model.RecipeKeys
	if err := c.BindJSON(&keys); err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidInput.Error())
		return
	}

	err = h.services.RecipeSharing.SetUserPublicKeyForRecipe(recipeId, userId, keys.PublicKey)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newMessageResponse(c, RespKeySet)
}

func (h *Handler) setRecipePrivateKey(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	var keys model.RecipeKeys
	if err := c.BindJSON(&keys); err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidInput.Error())
		return
	}

	err = h.services.RecipeSharing.SetUserPrivateKeyForRecipe(recipeId, userId, keys.PrivateKey)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newMessageResponse(c, RespKeySet)
}

func (h *Handler) deleteUserAccess(c *gin.Context) {
	requesterId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidInput.Error())
		return
	}

	err = h.services.RecipeSharing.DeleteUserAccessToRecipe(recipeId, userId, requesterId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newMessageResponse(c, RespKeyDeleted)
}