package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/models"
	"net/http"
	"strconv"
)

func (h *Handler) initCategoriesRoutes(api *gin.RouterGroup) {
	categories := api.Group("/categories", h.userIdentity)
	{
		categories.GET("", h.getCategories)
		categories.POST("/create", h.addCategory)
		categories.GET("/:category_id", h.getCategory)
		categories.PUT("/:category_id", h.updateCategory)
		categories.DELETE("/:category_id", h.deleteCategory)
	}
}

func (h *Handler) getCategories(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	categories, err := h.services.GetCategoriesByUser(userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, categories)
}

func (h *Handler) addCategory(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var category models.Category

	if err := c.BindJSON(&category); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	category.UserId = userId
	id, err := 	h.services.AddCategory(category)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"message": models.RespCategoryAdded,
	})
}

func (h *Handler) getCategory(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	categoryId, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	category, err := h.services.GetCategoryById(categoryId, userId)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if category.UserId != userId {
		newResponse(c, http.StatusForbidden, models.ErrAccessDenied.Error())
		return
	}
	c.JSON(http.StatusOK, category)
}

func (h *Handler) updateCategory(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	categoryId, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	var category models.Category

	if err := c.BindJSON(&category); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	category.Id = categoryId
	category.UserId = userId

	if err := h.services.UpdateCategory(category); err != nil {
		newResponse(c, http.StatusForbidden, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespCategoryUpdated,
	})
}

func (h *Handler) deleteCategory(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	categoryId, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	if err := h.services.DeleteCategory(categoryId, userId); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespCategoryDeleted,
	})
}