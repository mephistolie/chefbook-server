package v1

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) uploadUserKey(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	file, err := getFileByCtx(c)
	if err != nil {
		return
	}

	url, err := h.services.Encryption.UploadUserKey(c.Request.Context(), userId, file)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	newLinkResponse(c, url)
}

func (h *Handler) deleteUserKey(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	err = h.services.Encryption.DeleteUserKey(c.Request.Context(), userId)
	if err != nil {
		processKeyError(c, err)
		return
	}

	newMessageResponse(c, RespKeyDeleted)
}

func (h *Handler) getRecipeKey(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	url, err := h.services.Encryption.GetRecipeKey(recipeId, userId)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	newLinkResponse(c, url)
}

func (h *Handler) uploadRecipeKey(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	file, err := getFileByCtx(c)
	if err != nil {
		return
	}

	url, err := h.services.Encryption.UploadRecipeKey(c.Request.Context(), recipeId, userId, file)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	newLinkResponse(c, url)
}

func (h *Handler) deleteRecipeKey(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	err = h.services.Encryption.DeleteRecipeKey(c.Request.Context(), recipeId, userId)
	if err != nil {
		processKeyError(c, err)
		return
	}

	newMessageResponse(c, RespKeyDeleted)
}