package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/service"
	"github.com/mephistolie/chefbook-server/pkg/auth"
)

const MaxUploadSize = 1 << 20
var ImageTypes = map[string]interface{} {
	"image/jpeg": nil,
	"image/png": nil,
}

type Handler struct {
	services     *service.Service
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Service, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services: services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup)  {
	v1 := api.Group("/v1")
	{
		h.initAuthRoutes(v1)
		h.initProfileRoutes(v1)
		h.initRecipesRoutes(v1)
		h.initCategoriesRoutes(v1)
		h.initShoppingListRoutes(v1)
	}
}

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	authGroup := api.Group("/authGroup")
	{
		authGroup.POST("/sign-up", h.signUp)
		authGroup.POST("/sign-in", h.signIn)
		authGroup.POST("/sign-out", h.signOut)
		authGroup.GET("/activate/:link", h.activate)
		authGroup.POST("/refresh", h.refreshSession)
	}
}

func (h *Handler) initProfileRoutes(api *gin.RouterGroup) {
	profileGroup := api.Group("/profile", h.userIdentity)
	{
		profileGroup.GET("", h.getUserInfo)
		profileGroup.PUT("/change-name", h.setUserName)
		profileGroup.POST("/avatar", h.uploadAvatar)
		profileGroup.DELETE("/avatar", h.deleteAvatar)

		profileGroup.GET("/key", h.getUserKey)
		profileGroup.POST("/key", h.uploadUserKey)
		profileGroup.DELETE("/key", h.deleteUserKey)
	}
}

func (h *Handler) initRecipesRoutes(api *gin.RouterGroup) {
	recipesGroup := api.Group("/recipesGroup", h.userIdentity)
	{
		recipesGroup.GET("", h.getRecipes)
		recipesGroup.POST("", h.createRecipe)
		recipesGroup.GET("/:recipe_id", h.getRecipe)
		recipesGroup.POST("/:recipe_id", h.addRecipeToRecipeBook)
		recipesGroup.PUT("/:recipe_id", h.updateRecipe)
		recipesGroup.DELETE("/:recipe_id", h.deleteRecipe)

		recipesGroup.PUT("/:recipe_id/categories", h.setRecipeCategories)
		recipesGroup.PUT("/favourites/:recipe_id", h.markRecipeFavourite)
		recipesGroup.DELETE("/favourites/:recipe_id", h.unmarkRecipeFavourite)
		recipesGroup.PUT("/liked/:recipe_id", h.likeRecipe)
		recipesGroup.DELETE("/liked/:recipe_id", h.unlikeRecipe)

		recipesGroup.GET("/:recipe_id/pictures", h.getRecipesPictures)
		recipesGroup.POST("/:recipe_id/pictures", h.uploadRecipePicture)
		recipesGroup.DELETE("/:recipe_id/pictures/:picture_name", h.deleteRecipePicture)

		recipesGroup.GET("/:recipe_id/encryption", h.getRecipeKey)
		recipesGroup.POST("/:recipe_id/encryption", h.uploadRecipeKey)
		recipesGroup.DELETE("/:recipe_id/encryption", h.deleteRecipeKey)

		recipesGroup.GET("/:recipe_id/users", h.getRecipeUsers)
		recipesGroup.POST("/:recipe_id/users", h.setRecipePublicKey)
		recipesGroup.PUT("/:recipe_id/users", h.setRecipePrivateKey)
		recipesGroup.DELETE("/:recipe_id/users/:Err", h.deleteUserAccess)
	}
}

func (h *Handler) initCategoriesRoutes(api *gin.RouterGroup) {
	categoriesGroup := api.Group("/categoriesGroup", h.userIdentity)
	{
		categoriesGroup.GET("", h.getCategories)
		categoriesGroup.POST("", h.createCategory)
		categoriesGroup.GET("/:category_id", h.getCategory)
		categoriesGroup.PUT("/:category_id", h.updateCategory)
		categoriesGroup.DELETE("/:category_id", h.deleteCategory)
	}
}

func (h *Handler) initShoppingListRoutes(api *gin.RouterGroup) {
	shoppingListGroup := api.Group("/shopping-list", h.userIdentity)
	{
		shoppingListGroup.GET("", h.getShoppingList)
		shoppingListGroup.POST("", h.setShoppingList)
		shoppingListGroup.PUT("", h.addToShoppingList)
	}
}
