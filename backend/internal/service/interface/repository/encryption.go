package repository

import "github.com/google/uuid"

type Encryption interface {
	GetUserKeyLink(userId uuid.UUID) (*string, error)
	SetUserKeyLink(userId uuid.UUID, url *string) error
	GetRecipeKeyLink(recipeId uuid.UUID) (*string, error)
	SetRecipeKeyLink(recipeId uuid.UUID, url *string) error
}
