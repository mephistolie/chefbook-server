package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/app/dependencies/service"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/middleware"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/router/handler"
)

type v1Handler struct {
	auth            *handler.AuthHandler
	profile         *handler.ProfileHandler
	encryption      *handler.EncryptionHandler
	recipe          *handler.RecipeHandler
	recipeOwnership *handler.OwnedRecipeHandler
	recipePicture   *handler.RecipePictureHandler
	recipeSharing   *handler.RecipeSharingHandler
	category        *handler.CategoriesHandler
	shoppingList    *handler.ShoppingListHandler
}

type v1Router struct {
	middleware middleware.AuthMiddleware
	handler    v1Handler
}

func NewV1Router(services *service.Service, authMiddleware middleware.AuthMiddleware, fileMiddleware middleware.FileMiddleware) *v1Router {
	routesHandler := v1Handler{
		auth:            handler.NewAuthHandler(services.Auth),
		profile:         handler.NewProfileHandler(authMiddleware, fileMiddleware, services.Profile),
		encryption:      handler.NewEncryptionHandler(authMiddleware, fileMiddleware, services.Encryption),
		recipe:          handler.NewRecipeCrudHandler(authMiddleware, services.Recipe),
		recipeOwnership: handler.NewOwnedRecipeHandler(authMiddleware, services.RecipeOwnership),
		recipePicture:   handler.NewRecipePictureHandler(authMiddleware, fileMiddleware, services.RecipePicture),
		recipeSharing:   handler.NewRecipeSharingHandler(authMiddleware, fileMiddleware, services.RecipeSharing),
		category:        handler.NewCategoryHandler(authMiddleware, services.Category),
		shoppingList:    handler.NewShoppingListHandler(authMiddleware, services.ShoppingList),
	}

	return &v1Router{
		middleware: authMiddleware,
		handler:    routesHandler,
	}
}

func (r *v1Router) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		r.initAuthRoutes(v1)
		r.initProfileRoutes(v1)
		r.initRecipesRoutes(v1)
		r.initCategoriesRoutes(v1)
		r.initShoppingListRoutes(v1)
	}
}

func (r *v1Router) initAuthRoutes(api *gin.RouterGroup) {
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/sign-up", r.handler.auth.SignUp)
		authGroup.POST("/sign-in", r.handler.auth.SignIn)
		authGroup.POST("/sign-out", r.handler.auth.SignOut)
		authGroup.GET(fmt.Sprintf("/activate/:%s", handler.ParamActivationCode), r.handler.auth.ActivateProfile)
		authGroup.POST("/refresh", r.handler.auth.RefreshSession)
	}
}

func (r *v1Router) initProfileRoutes(api *gin.RouterGroup) {
	profileGroup := api.Group("/profile", r.middleware.CheckUserIdentity)
	{
		profileGroup.GET("", r.handler.profile.GetProfileInfo)
		profileGroup.PUT("/password", r.handler.profile.ChangePassword)
		profileGroup.PUT("/username", r.handler.profile.SetUsername)
		profileGroup.POST("/avatar", r.handler.profile.UploadAvatar)
		profileGroup.DELETE("/avatar", r.handler.profile.DeleteAvatar)

		profileGroup.GET("/key", r.handler.encryption.GetUserKey)
		profileGroup.POST("/key", r.handler.encryption.UploadUserKey)
		profileGroup.DELETE("/key", r.handler.encryption.DeleteUserKey)
	}
}

func (r *v1Router) initRecipesRoutes(api *gin.RouterGroup) {
	recipesGroup := api.Group("/recipes", r.middleware.CheckUserIdentity)
	{
		recipesGroup.GET("", r.handler.recipe.GetRecipes)
		recipesGroup.GET("/random", r.handler.recipe.GetRandomRecipe)

		recipesGroup.POST("", r.handler.recipeOwnership.CreateRecipe)
		recipesGroup.GET(fmt.Sprintf("/:%s", handler.ParamRecipeId), r.handler.recipe.GetRecipe)
		recipesGroup.PUT(fmt.Sprintf("/:%s", handler.ParamRecipeId), r.handler.recipeOwnership.UpdateRecipe)
		recipesGroup.DELETE(fmt.Sprintf("/:%s", handler.ParamRecipeId), r.handler.recipeOwnership.DeleteRecipe)

		recipesGroup.POST(fmt.Sprintf("/:%s/save", handler.ParamRecipeId), r.handler.recipe.AddRecipeToRecipeBook)
		recipesGroup.DELETE(fmt.Sprintf("/:%s/save", handler.ParamRecipeId), r.handler.recipe.RemoveFromRecipeBook)
		recipesGroup.PUT(fmt.Sprintf("/:%s/categories", handler.ParamRecipeId), r.handler.recipe.SetRecipeCategories)
		recipesGroup.PUT(fmt.Sprintf("/:%s/favourite", handler.ParamRecipeId), r.handler.recipe.MarkRecipeFavourite)
		recipesGroup.DELETE(fmt.Sprintf("/:%s/favourite", handler.ParamRecipeId), r.handler.recipe.UnmarkRecipeFavourite)
		recipesGroup.PUT(fmt.Sprintf("/:%s/likes", handler.ParamRecipeId), r.handler.recipe.LikeRecipe)
		recipesGroup.DELETE(fmt.Sprintf("/:%s/likes", handler.ParamRecipeId), r.handler.recipe.UnlikeRecipe)

		recipesGroup.GET(fmt.Sprintf("/:%s/pictures", handler.ParamRecipeId), r.handler.recipePicture.GetRecipePictures)
		recipesGroup.POST(fmt.Sprintf("/:%s/pictures", handler.ParamRecipeId), r.handler.recipePicture.UploadRecipePicture)
		recipesGroup.DELETE(fmt.Sprintf("/:%s/pictures/:%s", handler.ParamRecipeId, handler.ParamPictureId), r.handler.recipePicture.DeleteRecipePicture)

		recipesGroup.GET(fmt.Sprintf("/:%s/key", handler.ParamRecipeId), r.handler.encryption.GetRecipeKey)
		recipesGroup.POST(fmt.Sprintf("/:%s/key", handler.ParamRecipeId), r.handler.encryption.UploadRecipeKey)
		recipesGroup.DELETE(fmt.Sprintf("/:%s/key", handler.ParamRecipeId), r.handler.encryption.DeleteRecipeKey)

		recipesGroup.GET(fmt.Sprintf("/:%s/users", handler.ParamRecipeId), r.handler.recipeSharing.GetRecipeUsers)
		recipesGroup.GET(fmt.Sprintf("/:%s/users/key", handler.ParamRecipeId), r.handler.recipeSharing.GetUserRecipeKey)
		recipesGroup.PUT(fmt.Sprintf("/:%s/users/key", handler.ParamRecipeId), r.handler.recipeSharing.SetUserPublicKey)
		recipesGroup.DELETE(fmt.Sprintf("/:%s/users/key", handler.ParamRecipeId), r.handler.recipeSharing.DeleteUserPublicKey)
		recipesGroup.GET(fmt.Sprintf("/:%s/users/:%s/key", handler.ParamRecipeId, handler.ParamUserId), r.handler.recipeSharing.GetUserPublicKey)
		recipesGroup.PUT(fmt.Sprintf("/:%s/users/:%s/key", handler.ParamRecipeId, handler.ParamUserId), r.handler.recipeSharing.SetOwnerPrivateKey)
		recipesGroup.DELETE(fmt.Sprintf("/:%s/users/:%s/key", handler.ParamRecipeId, handler.ParamUserId), r.handler.recipeSharing.DeleteOwnerPrivateKey)
		recipesGroup.DELETE(fmt.Sprintf("/:%s/users/:%s", handler.ParamRecipeId, handler.ParamUserId), r.handler.recipeSharing.DeleteUserAccess)
	}
}

func (r *v1Router) initCategoriesRoutes(api *gin.RouterGroup) {
	categoriesGroup := api.Group("/categories", r.middleware.CheckUserIdentity)
	{
		categoriesGroup.GET("", r.handler.category.GetCategories)
		categoriesGroup.POST("", r.handler.category.CreateCategory)
		categoriesGroup.GET(fmt.Sprintf("/:%s", handler.ParamCategoryId), r.handler.category.GetCategory)
		categoriesGroup.PUT(fmt.Sprintf("/:%s", handler.ParamCategoryId), r.handler.category.UpdateCategory)
		categoriesGroup.DELETE(fmt.Sprintf("/:%s", handler.ParamCategoryId), r.handler.category.DeleteCategory)
	}
}

func (r *v1Router) initShoppingListRoutes(api *gin.RouterGroup) {
	shoppingListGroup := api.Group("/shopping-list", r.middleware.CheckUserIdentity)
	{
		shoppingListGroup.GET("", r.handler.shoppingList.GetShoppingList)
		shoppingListGroup.POST("", r.handler.shoppingList.SetShoppingList)
		shoppingListGroup.PUT("", r.handler.shoppingList.AddToShoppingList)
	}
}
