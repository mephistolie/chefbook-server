package repository

import "github.com/mephistolie/chefbook-server/internal/entity"

type RecipeOwnership interface {
	CreateRecipe(recipe entity.RecipeInput, userId int) (int, error)
	UpdateRecipe(recipeId int, recipe entity.RecipeInput) error
	DeleteRecipe(recipeId int) error
}

type Recipe interface {
	GetRecipes(params entity.RecipesQuery, userId int) ([]entity.RecipeInfo, error)
	GetRecipe(recipeId int) (entity.Recipe, error)
	GetRandomRecipe(languages *[]string, userId int) (entity.UserRecipe, error)
	GetRecipeWithUserFields(recipeId int, userId int) (entity.UserRecipe, error)
	GetRecipeOwnerId(recipeId int) (int, error)
	AddRecipeToRecipeBook(recipeId, userId int) error
	RemoveRecipeFromRecipeBook(recipeId, userId int) error
	SetRecipeCategories(recipeId int, categoriesIds []int, userId int) error
	SetRecipeFavourite(recipeId int, isFavourite bool, userId int) error
	SetRecipeLiked(recipeId int, isLiked bool, userId int) error
}

type RecipeSharing interface {
	GetRecipeUserList(recipeId int) ([]entity.ProfileInfo, error)
	GetUserPublicKey(recipeId, userId int) (string, error)
	SetUserPublicKeyLink(recipeId int, userId int, userKey *string) error
	GetUserRecipeKey(recipeId, userId int) (string, error)
	SetOwnerPrivateKeyLinkForUser(recipeId int, userId int, userKey *string) error
}
