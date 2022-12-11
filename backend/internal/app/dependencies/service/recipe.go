package service

import (
	"context"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Recipe interface {
	GetRecipes(query entity.RecipesQuery, userId string) ([]entity.RecipeInfo, error)
	GetRecipe(recipeId, userId string) (entity.UserRecipe, error)
	GetRandomRecipe(languages *[]string, userId string) (entity.UserRecipe, error)
	AddRecipeToRecipeBook(recipeId, userId string) error
	RemoveRecipeFromRecipeBook(recipeId, userId string) error
	SetRecipeCategories(recipeId string, categories []string, userId string) error
	SetRecipeFavourite(recipeId string, favourite bool, userId string) error
	SetRecipeLikeStatus(recipeId string, favourite bool, userId string) error
}

type RecipeOwnership interface {
	CreateRecipe(recipe entity.RecipeInput, userId string) (string, error)
	UpdateRecipe(recipe entity.RecipeInput, recipeId, userId string) error
	DeleteRecipe(recipeId, userId string) error
}

type RecipePicture interface {
	GetRecipePictures(ctx context.Context, recipeId string, userId string) ([]string, error)
	UploadRecipePicture(ctx context.Context, recipeId, userId string, file entity.MultipartFile) (string, error)
	DeleteRecipePicture(ctx context.Context, recipeId, userId string, pictureName string) error
}

type RecipeSharing interface {
	GetUsersList(recipeId, userId string) ([]entity.ProfileInfo, error)
	GetUserPublicKey(recipeId, userId, requesterId string) (string, error)
	SetUserPublicKey(recipeId string, userId string, userKey *string) error
	GetOwnerPrivateKeyForUser(recipeId, userId string) (string, error)
	SetOwnerPrivateKeyForUser(recipeId string, userId string, requesterId string, ownerKey *string) error
	DeleteUserAccess(recipeId, userId, requesterId string) error
}
