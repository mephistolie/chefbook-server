package handler

import (
	"chefbook-server/internal/app/dependencies/service"
	"chefbook-server/internal/delivery/http/middleware"
	"chefbook-server/internal/delivery/http/middleware/response"
	"chefbook-server/internal/delivery/http/presentation/response_body/message"
	"github.com/gin-gonic/gin"
)

const (
	ParamPictureId = "picture_id"

	maxRecipePictureSize = 1 << 20
)

type RecipePictureHandler struct {
	authMiddleware middleware.AuthMiddleware
	fileMiddleware middleware.FileMiddleware
	service        service.RecipePicture
}

func NewRecipePictureHandler(authMiddleware middleware.AuthMiddleware, fileMiddleware middleware.FileMiddleware, service service.RecipePicture) *RecipePictureHandler {
	return &RecipePictureHandler{
		authMiddleware: authMiddleware,
		fileMiddleware: fileMiddleware,
		service:        service,
	}
}

// GetRecipePictures Swagger Documentation
// @Summary Get Recipe Pictures
// @Security ApiKeyAuth
// @Tags recipe-pictures
// @Description Get recipe pictures links
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} []string
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/pictures [get]
func (r *RecipePictureHandler) GetRecipePictures(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	pictures, err := r.service.GetRecipePictures(c.Request.Context(), recipeId, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Success(c, pictures)
}

// UploadRecipePicture Swagger Documentation
// @Summary Upload Recipe Picture
// @Security ApiKeyAuth
// @Tags recipe-pictures
// @Description Upload recipe picture as usual or encrypted file
// @Accept mpfd
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Param file formData file true "Picture File"
// @Success 200 {object} response_body.Link
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/pictures [post]
func (r *RecipePictureHandler) UploadRecipePicture(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	file, err := r.fileMiddleware.GetFileWithMaxSize(c, maxRecipePictureSize)
	if err != nil {
		response.Failure(c, err)
		return
	}

	url, err := r.service.UploadRecipePicture(c.Request.Context(), recipeId, userId, file)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Link(c, url)
}

// DeleteRecipePicture Swagger Documentation
// @Summary Delete Recipe Picture
// @Security ApiKeyAuth
// @Tags recipe-pictures
// @Description Delete recipe picture by name
// @Accept mpfd
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Param picture_name path int true "Picture Name"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/pictures/{picture_name} [delete]
func (r *RecipePictureHandler) DeleteRecipePicture(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	pictureName := c.Param(ParamPictureId)

	err = r.service.DeleteRecipePicture(c.Request.Context(), recipeId, userId, pictureName)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.RecipePictureDeleted)
}
