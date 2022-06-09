package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/app/dependencies/service"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/middleware"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/middleware/response"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/response_body/message"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"strconv"
)

const (
	maxKeySize = 1 << 20
)

type EncryptionHandler struct {
	authMiddleware middleware.AuthMiddleware
	fileMiddleware middleware.FileMiddleware
	service        service.Encryption
}

func NewEncryptionHandler(authMiddleware middleware.AuthMiddleware, fileMiddleware middleware.FileMiddleware, service service.Encryption) *EncryptionHandler {
	return &EncryptionHandler{
		authMiddleware: authMiddleware,
		fileMiddleware: fileMiddleware,
		service:        service,
	}
}

// GetUserKey Swagger Documentation
// @Summary Get User Key
// @Security ApiKeyAuth
// @Tags profile-encryption
// @Description Get user encrypted vault key (AES encrypted by generated RSA)
// @Accept json
// @Produce json
// @Success 200 {object} response_body.Link
// @Failure 400 {object} response_body.Error
// @Router /v1/profile/key [get]
func (r *EncryptionHandler) GetUserKey(c *gin.Context) {
	userId, err := r.authMiddleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	url, err := r.service.GetUserKeyLink(userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Link(c, url)
}

// UploadUserKey Swagger Documentation
// @Summary Upload User Key
// @Security ApiKeyAuth
// @Tags profile-encryption
// @Description Upload user encrypted vault key (RSA encrypted by vault password generated AES)
// @Accept mpfd
// @Produce json
// @Param file formData file true "User Key File"
// @Success 200 {object} response_body.Link
// @Failure 400 {object} response_body.Error
// @Router /v1/profile/key [post]
func (r *EncryptionHandler) UploadUserKey(c *gin.Context) {
	userId, err := r.authMiddleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	file, err := r.fileMiddleware.GetFileWithMaxSize(c, maxKeySize)
	if err != nil {
		response.Failure(c, err)
		return
	}

	url, err := r.service.UploadUserKey(c.Request.Context(), userId, file)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Link(c, url)
}

// DeleteUserKey Swagger Documentation
// @Summary Delete User Key
// @Security ApiKeyAuth
// @Tags profile-encryption
// @Description Delete user encrypted vault key (RSA encrypted by vault password generated AES)
// @Accept json
// @Produce json
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/profile/key [delete]
func (r *EncryptionHandler) DeleteUserKey(c *gin.Context) {
	userId, err := r.authMiddleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	err = r.service.DeleteUserKey(c.Request.Context(), userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.KeyDeleted)
}

// GetRecipeKey Swagger Documentation
// @Summary Get Recipe Key
// @Security ApiKeyAuth
// @Tags recipe-encryption
// @Description Get recipe encrypted key (AES encrypted by user RSA Private Key)
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} response_body.Link
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/key [get]
func (r *EncryptionHandler) GetRecipeKey(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	url, err := r.service.GetRecipeKey(recipeId, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Link(c, url)
}

// UploadRecipeKey Swagger Documentation
// @Summary Upload Recipe Key
// @Security ApiKeyAuth
// @Tags recipe-encryption
// @Description Upload recipe encrypted key (AES encrypted by user RSA Private Key)
// @Accept mpfd
// @Produce json
// @Param recipe_id path string true "Recipe ID"
// @Param file formData file true "Recipe Key File"
// @Success 200 {object} response_body.Link
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/key [post]
func (r *EncryptionHandler) UploadRecipeKey(c *gin.Context) {
	userId, err := r.authMiddleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	recipeId, err := strconv.Atoi(c.Param(ParamRecipeId))
	if err != nil {
		response.Failure(c, failure.Unknown)
		return
	}

	file, err := r.fileMiddleware.GetFileWithMaxSize(c, maxKeySize)
	if err != nil {
		response.Failure(c, err)
		return
	}

	url, err := r.service.UploadRecipeKey(c.Request.Context(), recipeId, userId, file)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Link(c, url)
}

// DeleteRecipeKey Swagger Documentation
// @Summary Delete Recipe Key
// @Security ApiKeyAuth
// @Tags recipe-encryption
// @Description Delete recipe encrypted key (AES encrypted by user RSA Private Key)
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/key [delete]
func (r *EncryptionHandler) DeleteRecipeKey(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		return
	}

	err = r.service.DeleteRecipeKey(c.Request.Context(), recipeId, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.KeyDeleted)
}
