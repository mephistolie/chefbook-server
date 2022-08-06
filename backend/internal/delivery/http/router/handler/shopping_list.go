package handler

import (
	"chefbook-server/internal/app/dependencies/service"
	"chefbook-server/internal/delivery/http/middleware"
	"chefbook-server/internal/delivery/http/middleware/response"
	"chefbook-server/internal/delivery/http/presentation/common_body"
	"chefbook-server/internal/delivery/http/presentation/request_body"
	"chefbook-server/internal/delivery/http/presentation/response_body"
	"chefbook-server/internal/delivery/http/presentation/response_body/message"
	"chefbook-server/internal/entity/failure"
	"github.com/gin-gonic/gin"
)

type ShoppingListHandler struct {
	middleware middleware.AuthMiddleware
	service    service.ShoppingList
}

func NewShoppingListHandler(middleware middleware.AuthMiddleware, service service.ShoppingList) *ShoppingListHandler {
	return &ShoppingListHandler{
		middleware: middleware,
		service:    service,
	}
}

// GetShoppingList Swagger Documentation
// @Summary Get Shopping List
// @Security ApiKeyAuth
// @Tags shopping-list
// @Description Get user shopping list
// @Accept json
// @Produce json
// @Success 200 {object} response_body.ShoppingList
// @Failure 400 {object} response_body.Error
// @Router /v1/shopping-list [get]
func (r *ShoppingListHandler) GetShoppingList(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	shoppingList, err := r.service.GetShoppingList(userId)
	if err != nil {
		response.Failure(c, err)
	}

	response.Success(c, response_body.NewShoppingList(shoppingList))
}

// SetShoppingList Swagger Documentation
// @Summary Set Shopping List
// @Security ApiKeyAuth
// @Tags shopping-list
// @Description Set user shopping list
// @Accept json
// @Produce json
// @Param input body []common_body.Purchase true "Shopping List"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/shopping-list [post]
func (r *ShoppingListHandler) SetShoppingList(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	var shoppingList []common_body.Purchase
	if err := c.BindJSON(&shoppingList); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	if err := r.service.SetShoppingList(request_body.NewShoppingListEntity(shoppingList), userId); err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.ShoppingListUpdated)
}

// AddToShoppingList Swagger Documentation
// @Summary Add Purchases to Shopping List
// @Security ApiKeyAuth
// @Tags shopping-list
// @Description Add purchases to user shopping list
// @Accept json
// @Produce json
// @Param input body []common_body.Purchase true "New Purchases"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/shopping-list [put]
func (r *ShoppingListHandler) AddToShoppingList(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	var purchases []common_body.Purchase
	if err := c.BindJSON(&purchases); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	if err := r.service.AddToShoppingList(request_body.NewShoppingListEntity(purchases), userId); err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.ShoppingListUpdated)
}
