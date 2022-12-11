package repository

import "github.com/mephistolie/chefbook-server/internal/entity"

type RecipeOwnership interface {
	CreateRecipe(recipe entity.RecipeInput, userId string) (string, error)
	UpdateRecipe(recipeId string, recipe entity.RecipeInput) error
	DeleteRecipe(recipeId string) error
}

type Recipe interface {
	GetRecipes(params entity.RecipesQuery, userId string) ([]entity.RecipeInfo, error)
	GetRecipe(recipeId string) (entity.Recipe, error)
	GetRandomRecipe(languages *[]string, userId string) (entity.UserRecipe, error)
	GetRecipeWithUserFields(recipeId, userId string) (entity.UserRecipe, error)
	GetRecipeOwnerId(recipeId string) (string, error)
	AddRecipeToRecipeBook(recipeId, userId string) error
	RemoveRecipeFromRecipeBook(recipeId, userId string) error
	SetRecipeCategories(recipeId string, categoriesIds []string, userId string) error
	SetRecipeFavourite(recipeId string, isFavourite bool, userId string) error
	SetRecipeLiked(recipeId string, isLiked bool, userId string) error
}

type RecipeSharing interface {
	GetRecipeUserList(recipeId string) ([]entity.ProfileInfo, error)
	GetUserPublicKey(recipeId, userId string) (string, error)
	SetUserPublicKeyLink(recipeId string, userId string, userKey *string) error
	GetUserRecipeKey(recipeId, userId string) (string, error)
	SetOwnerPrivateKeyLinkForUser(recipeId string, userId string, userKey *string) error
}
