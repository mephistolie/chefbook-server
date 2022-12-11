package repository

type Encryption interface {
	GetUserKeyLink(userId string) (*string, error)
	SetUserKeyLink(userId string, url *string) error
	GetRecipeKeyLink(recipeId string) (*string, error)
	SetRecipeKeyLink(recipeId string, url *string) error
}
