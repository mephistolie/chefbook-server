package service

import (
	"chefbook-server/internal/entity"
	"context"
)

type Recipe interface {
	GetRecipes(query entity.RecipesQuery, userId int) ([]entity.RecipeInfo, error)
	GetRecipe(recipeId, userId int) (entity.UserRecipe, error)
	GetRandomRecipe(languages *[]string, userId int) (entity.UserRecipe, error)
	AddRecipeToRecipeBook(recipeId, userId int) error
	RemoveRecipeFromRecipeBook(recipeId, userId int) error
	SetRecipeCategories(recipeId int, categories []int, userId int) error
	SetRecipeFavourite(recipeId int, favourite bool, userId int) error
	SetRecipeLikeStatus(recipeId int, favourite bool, userId int) error
}

type RecipeOwnership interface {
	CreateRecipe(recipe entity.RecipeInput, userId int) (int, error)
	UpdateRecipe(recipe entity.RecipeInput, recipeId, userId int) error
	DeleteRecipe(recipeId, userId int) error
}

type RecipePicture interface {
	GetRecipePictures(ctx context.Context, recipeId int, userId int) ([]string, error)
	UploadRecipePicture(ctx context.Context, recipeId, userId int, file entity.MultipartFile) (string, error)
	DeleteRecipePicture(ctx context.Context, recipeId, userId int, pictureName string) error
}

type RecipeSharing interface {
	GetUsersList(recipeId, userId int) ([]entity.ProfileInfo, error)
	GetUserPublicKey(recipeId, userId, requesterId int) (string, error)
	SetUserPublicKey(recipeId int, userId int, userKey *string) error
	GetOwnerPrivateKeyForUser(recipeId, userId int) (string, error)
	SetOwnerPrivateKeyForUser(recipeId int, userId int, requesterId int, ownerKey *string) error
	DeleteUserAccess(recipeId, userId, requesterId int) error
}
