package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/model"
	"net/http"
)

func (h *Handler) getRecipesPictures(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	objects, err := h.services.RecipePictures.GetRecipePictures(c.Request.Context(), recipeId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, model.ErrUnableDeleteRecipePicture.Error())
		return
	}

	c.JSON(http.StatusOK, objects)
}


func (h *Handler) uploadRecipePicture(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}

	file, err := getFileByCtx(c)
	if err != nil {
		return
	}

	url, err := h.services.RecipePictures.UploadRecipePicture(c.Request.Context(), recipeId, userId, file)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newLinkResponse(c, url)
}

func (h *Handler) deleteRecipePicture(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIdByCtx(c)
	if err != nil {
		return
	}
	pictureName := c.Param("picture_name")

	err = h.services.RecipePictures.DeleteRecipePicture(c.Request.Context(), recipeId, userId, pictureName)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, model.ErrUnableDeleteRecipePicture.Error())
		return
	}

	newMessageResponse(c, RespRecipePictureDeleted)
}