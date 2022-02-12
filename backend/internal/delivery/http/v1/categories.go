package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getCategories(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	categories, err := h.services.Categories.GetUserCategories(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, categories)
}

func (h *Handler) createCategory(c *gin.Context) {
	category, err := getCategoryByCtx(c)
	if err != nil {
		return
	}

	id, err := 	h.services.Categories.AddCategory(category)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newIdResponse(c, id, RespCategoryAdded)
}

func (h *Handler) getCategory(c *gin.Context) {
	userId, categoryId, err := getUserAndCategoryIdByCtx(c)
	if err != nil {
		return
	}

	category, err := h.services.Categories.GetCategoryById(categoryId, userId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *Handler) updateCategory(c *gin.Context) {
	category, err := getCategoryByCtx(c)
	if err != nil {
		return
	}

	if err := h.services.Categories.UpdateCategory(category); err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	newMessageResponse(c, RespCategoryUpdated)
}

func (h *Handler) deleteCategory(c *gin.Context) {
	userId, categoryId, err := getUserAndCategoryIdByCtx(c)
	if err != nil {
		return
	}

	if err := h.services.Categories.DeleteCategory(categoryId, userId); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newMessageResponse(c, RespCategoryDeleted)
}