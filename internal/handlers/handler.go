package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/services"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler  {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine  {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api")
	{
		recipes := api.Group("/recipes")
		{
			recipes.GET("/", h.readRecipe)
			recipes.POST("/", h.createRecipe)
			recipes.GET("/:id", h.readRecipe)
			recipes.PUT("/:id", h.updateRecipe)
			recipes.DELETE("/:id", h.deleteRecipe)
		}
	}

	return router
}