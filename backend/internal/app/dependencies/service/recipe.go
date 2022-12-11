package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
)

type Recipe interface {
	GetRecipes(query entity.RecipesQuery, userId uuid.UUID) ([]entity.RecipeInfo, error)
	GetRecipe(recipeId, userId uuid.UUID) (entity.UserRecipe, error)
	GetRandomRecipe(languages *[]string, userId uuid.UUID) (entity.UserRecipe, error)
	AddRecipeToRecipeBook(recipeId, userId uuid.UUID) error
	RemoveRecipeFromRecipeBook(recipeId, userId uuid.UUID) error
	SetRecipeCategories(recipeId uuid.UUID, categories []uuid.UUID, userId uuid.UUID) error
	SetRecipeFavourite(recipeId uuid.UUID, favourite bool, userId uuid.UUID) error
	SetRecipeLikeStatus(recipeId uuid.UUID, favourite bool, userId uuid.UUID) error
}

type RecipeOwnership interface {
	CreateRecipe(recipe entity.RecipeInput, userId uuid.UUID) (uuid.UUID, error)
	UpdateRecipe(recipe entity.RecipeInput, recipeId, userId uuid.UUID) error
	DeleteRecipe(recipeId, userId uuid.UUID) error
}

type RecipePicture interface {
	GetRecipePictures(ctx context.Context, recipeId uuid.UUID, userId uuid.UUID) ([]string, error)
	UploadRecipePicture(ctx context.Context, recipeId, userId uuid.UUID, file entity.MultipartFile) (string, error)
	DeleteRecipePicture(ctx context.Context, recipeId, userId uuid.UUID, pictureName string) error
}

type RecipeSharing interface {
	GetUsersList(recipeId, userId uuid.UUID) ([]entity.ProfileInfo, error)
	GetUserPublicKey(recipeId, userId, requesterId uuid.UUID) (string, error)
	SetUserPublicKey(recipeId uuid.UUID, userId uuid.UUID, userKey *string) error
	GetOwnerPrivateKeyForUser(recipeId, userId uuid.UUID) (string, error)
	SetOwnerPrivateKeyForUser(recipeId uuid.UUID, userId uuid.UUID, requesterId uuid.UUID, ownerKey *string) error
	DeleteUserAccess(recipeId, userId, requesterId uuid.UUID) error
}
