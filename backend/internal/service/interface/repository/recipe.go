package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type RecipeOwnership interface {
	CreateRecipe(recipe entity.RecipeInput, userId uuid.UUID) (uuid.UUID, error)
	UpdateRecipe(recipeId uuid.UUID, recipe entity.RecipeInput) error
	DeleteRecipe(recipeId uuid.UUID) error
}

type Recipe interface {
	GetRecipes(params entity.RecipesQuery, userId uuid.UUID) ([]entity.RecipeInfo, error)
	GetRecipe(recipeId uuid.UUID) (entity.Recipe, error)
	GetRandomRecipe(languages *[]string, userId uuid.UUID) (entity.UserRecipe, error)
	GetRecipeWithUserFields(recipeId, userId uuid.UUID) (entity.UserRecipe, error)
	GetRecipeOwnerId(recipeId uuid.UUID) (uuid.UUID, error)
	AddRecipeToRecipeBook(recipeId, userId uuid.UUID) error
	RemoveRecipeFromRecipeBook(recipeId, userId uuid.UUID) error
	SetRecipeCategories(recipeId uuid.UUID, categoriesIds []uuid.UUID, userId uuid.UUID) error
	SetRecipeFavourite(recipeId uuid.UUID, isFavourite bool, userId uuid.UUID) error
	SetRecipeLiked(recipeId uuid.UUID, isLiked bool, userId uuid.UUID) error
}

type RecipeSharing interface {
	GetRecipeUserList(recipeId uuid.UUID) ([]entity.ProfileInfo, error)
	GetUserPublicKey(recipeId, userId uuid.UUID) (string, error)
	SetUserPublicKeyLink(recipeId uuid.UUID, userId uuid.UUID, userKey *string) error
	GetUserRecipeKey(recipeId, userId uuid.UUID) (string, error)
	SetOwnerPrivateKeyLinkForUser(recipeId uuid.UUID, userId uuid.UUID, userKey *string) error
}
