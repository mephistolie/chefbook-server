package handler

import (
	"chefbook-server/internal/app/dependencies/service"
	"chefbook-server/internal/delivery/http/middleware"
	"chefbook-server/internal/delivery/http/middleware/response"
	"chefbook-server/internal/delivery/http/presentation/request_body"
	"chefbook-server/internal/delivery/http/presentation/response_body"
	"chefbook-server/internal/delivery/http/presentation/response_body/message"
	"chefbook-server/internal/entity/failure"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	ParamCategoryId = "category_id"
)

type CategoriesHandler struct {
	middleware middleware.AuthMiddleware
	service    service.Category
}

func NewCategoryHandler(middleware middleware.AuthMiddleware, service service.Category) *CategoriesHandler {
	return &CategoriesHandler{
		middleware: middleware,
		service:    service,
	}
}

// GetCategories Swagger Documentation
// @Summary Get Categories
// @Security ApiKeyAuth
// @Tags categories
// @Description Get user categories
// @Accept json
// @Produce json
// @Success 200 {object} []response_body.Category
// @Failure 400 {object} response_body.Error
// @Router /v1/categories [get]
func (r *CategoriesHandler) GetCategories(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	categories := r.service.GetUserCategories(userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Success(c, response_body.NewCategories(categories))
}

// CreateCategory Swagger Documentation
// @Summary Create Category
// @Security ApiKeyAuth
// @Tags categories
// @Description Create new user category
// @Accept json
// @Produce json
// @Param input body request_body.CategoryInput true "Category"
// @Success 200 {object} response_body.Id
// @Failure 400 {object} response_body.Error
// @Router /v1/categories [post]
func (r *CategoriesHandler) CreateCategory(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	var body request_body.CategoryInput
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	categoryId, err := r.service.CreateCategory(body.Entity(), userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.NewId(c, categoryId, message.CategoryCreated)
}

// GetCategory Swagger Documentation
// @Summary Get Category
// @Security ApiKeyAuth
// @Tags categories
// @Description Get user category
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {object} response_body.Category
// @Failure 400 {object} response_body.Error
// @Router /v1/categories/{category_id} [get]
func (r *CategoriesHandler) GetCategory(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	categoryId, err := strconv.Atoi(c.Param(ParamCategoryId))
	if err != nil {
		response.Failure(c, failure.Unknown)
		return
	}

	category, err := r.service.GetCategory(categoryId, userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Success(c, response_body.NewCategory(category))
}

// UpdateCategory Swagger Documentation
// @Summary Update Category
// @Security ApiKeyAuth
// @Tags categories
// @Description Update user category
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Param input body request_body.CategoryInput true "Category"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/categories/{category_id} [put]
func (r *CategoriesHandler) UpdateCategory(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	categoryId, err := strconv.Atoi(c.Param(ParamCategoryId))
	if err != nil {
		response.Failure(c, failure.Unknown)
		return
	}

	var body request_body.CategoryInput
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	if err := r.service.UpdateCategory(categoryId, body.Entity(), userId); err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.CategoryUpdated)
}

// DeleteCategory Swagger Documentation
// @Summary Delete Category
// @Security ApiKeyAuth
// @Tags categories
// @Description Delete user category
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/categories/{category_id} [delete]
func (r *CategoriesHandler) DeleteCategory(c *gin.Context) {
	userId, err := r.middleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	categoryId, err := strconv.Atoi(c.Param(ParamCategoryId))
	if err != nil {
		response.Failure(c, failure.Unknown)
		return
	}

	if err := r.service.DeleteCategory(categoryId, userId); err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.CategoryDeleted)
}
