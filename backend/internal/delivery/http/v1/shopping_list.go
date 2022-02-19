package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/model"
	"net/http"
)
func (h *Handler) getShoppingList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	shoppingList, err := h.services.ShoppingList.GetShoppingList(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, shoppingList)
}

func (h *Handler) addToShoppingList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var purchases []model.Purchase
	if err := c.BindJSON(&purchases); err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidInput.Error())
		return
	}

	if err := 	h.services.ShoppingList.AddToShoppingList(purchases, userId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newMessageResponse(c, RespShoppingListUpdated)
}

func (h *Handler) setShoppingList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var shoppingList model.ShoppingList
	if err := c.BindJSON(&shoppingList); err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidInput.Error())
		return
	}

	if err := 	h.services.ShoppingList.SetShoppingList(shoppingList, userId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newMessageResponse(c, RespShoppingListUpdated)
}