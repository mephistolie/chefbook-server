package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/app/dependencies/service"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/middleware"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/middleware/response"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/common_body"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/response_body"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/response_body/message"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
)

type RecipeSharingHandler struct {
	authMiddleware middleware.AuthMiddleware
	fileMiddleware middleware.FileMiddleware
	service        service.RecipeSharing
}

func NewRecipeSharingHandler(authMiddleware middleware.AuthMiddleware, fileMiddleware middleware.FileMiddleware, service service.RecipeSharing) *RecipeSharingHandler {
	return &RecipeSharingHandler{
		authMiddleware: authMiddleware,
		fileMiddleware: fileMiddleware,
		service:        service,
	}
}

// GetRecipeUsers Swagger Documentation
// @Summary Get Recipe Users
// @Security ApiKeyAuth
// @Tags recipe-sharing
// @Description Get users that saved recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} []response_body.MinimalProfileInfo
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/users [get]
func (r *RecipeSharingHandler) GetRecipeUsers(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	userList, err := r.service.GetUsersList(recipeId, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Success(c, response_body.NewUsersList(userList))
}

// GetUserPublicKey Swagger Documentation
// @Summary Get User Public Key for Recipe
// @Security ApiKeyAuth
// @Tags recipe-sharing
// @Description Get user public key for recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} common_body.RecipeUserPublicKey
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/users/{user_id}/key [get]
func (r *RecipeSharingHandler) GetUserPublicKey(c *gin.Context) {
	requesterId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	userId, err := uuid.Parse(c.Param(ParamRecipeId))
	if err != nil {
		response.Failure(c, err)
		return
	}

	var body common_body.RecipeUserPublicKey
	body.PublicKey, err = r.service.GetUserPublicKey(recipeId, userId, requesterId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Success(c, body)
}

// SetUserPublicKey Swagger Documentation
// @Summary Set Recipe User Public Key
// @Security ApiKeyAuth
// @Tags recipe-sharing
// @Description Set user public key for recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Param input body common_body.RecipeUserPublicKey true "Key"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/users/key [put]
func (r *RecipeSharingHandler) SetUserPublicKey(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	var body common_body.RecipeUserPublicKey
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	err = r.service.SetUserPublicKey(recipeId, userId, &body.PublicKey)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.KeySet)
}

// DeleteUserPublicKey Swagger Documentation
// @Summary Delete Recipe User Public Key
// @Security ApiKeyAuth
// @Tags recipe-sharing
// @Description Delete user public key for recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/users/key [delete]
func (r *RecipeSharingHandler) DeleteUserPublicKey(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	err = r.service.SetUserPublicKey(recipeId, userId, nil)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.KeyDeleted)
}

// GetUserRecipeKey Swagger Documentation
// @Summary Get User Recipe Private Key
// @Security ApiKeyAuth
// @Tags recipe-sharing
// @Description Get recipe key encrypted by user public key
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Success 200 {object} common_body.RecipeOwnerPrivateKey
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/users/key [get]
func (r *RecipeSharingHandler) GetUserRecipeKey(c *gin.Context) {
	userId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	var body common_body.RecipeOwnerPrivateKey
	body.PrivateKey, err = r.service.GetOwnerPrivateKeyForUser(recipeId, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Success(c, body)
}

// SetOwnerPrivateKey Swagger Documentation
// @Summary Set Recipe Owner Private Key
// @Security ApiKeyAuth
// @Tags recipe-sharing
// @Description Set owner private key encrypted by public user key for recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Param user_id path int true "User ID"
// @Param input body common_body.RecipeOwnerPrivateKey true "Key"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/users/{user_id}/key [post]
func (r *RecipeSharingHandler) SetOwnerPrivateKey(c *gin.Context) {
	requesterId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	userId, err := uuid.Parse(c.Param(ParamRecipeId))
	if err != nil {
		response.Failure(c, err)
		return
	}

	var body common_body.RecipeOwnerPrivateKey
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	err = r.service.SetOwnerPrivateKeyForUser(recipeId, userId, requesterId, &body.PrivateKey)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.KeySet)
}

// DeleteOwnerPrivateKey Swagger Documentation
// @Summary Delete Recipe Owner Private Key
// @Security ApiKeyAuth
// @Tags recipe-sharing
// @Description Delete owner private key encrypted by public user key for recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/users/{user_id}/key [delete]
func (r *RecipeSharingHandler) DeleteOwnerPrivateKey(c *gin.Context) {
	requesterId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	userId, err := uuid.Parse(c.Param(ParamRecipeId))
	if err != nil {
		response.Failure(c, err)
		return
	}

	err = r.service.SetOwnerPrivateKeyForUser(recipeId, userId, requesterId, nil)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.KeyDeleted)
}

// DeleteUserAccess Swagger Documentation
// @Summary Delete User Access to Recipe
// @Security ApiKeyAuth
// @Tags recipe-sharing
// @Description Delete user access to recipe
// @Accept json
// @Produce json
// @Param recipe_id path int true "Recipe ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/recipes/{recipe_id}/users/{user_id} [delete]
func (r *RecipeSharingHandler) DeleteUserAccess(c *gin.Context) {
	requesterId, recipeId, err := getUserAndRecipeIds(c, r.authMiddleware)
	if err != nil {
		response.Failure(c, err)
		return
	}

	userId, err := uuid.Parse(c.Param(ParamRecipeId))
	if err != nil {
		response.Failure(c, err)
		return
	}

	err = r.service.DeleteUserAccess(recipeId, userId, requesterId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.KeyDeleted)
}
