package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/models"
	"net/http"
)

func (h *Handler) initShoppingListRoutes(api *gin.RouterGroup) {
	shoppingList := api.Group("/shopping-list", h.userIdentity)
	{
		shoppingList.GET("", h.getShoppingList)
		shoppingList.POST("", h.setShoppingList)
		shoppingList.PUT("", h.addToShoppingList)
	}
}

func (h *Handler) getShoppingList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shoppingList, err := h.services.ShoppingList.GetShoppingList(userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, shoppingList)
}

func (h *Handler) addToShoppingList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var purchases []models.Purchase
	if err := c.BindJSON(&purchases); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	if err := 	h.services.ShoppingList.AddToShoppingList(purchases, userId); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespShoppingListUpdated,
	})
}

func (h *Handler) setShoppingList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var shoppingList models.ShoppingList
	if err := c.BindJSON(&shoppingList); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	if err := 	h.services.ShoppingList.SetShoppingList(shoppingList, userId); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespShoppingListUpdated,
	})
}