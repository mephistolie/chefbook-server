package repository

type Encryption interface {
	GetUserKeyLink(userId int) (*string, error)
	SetUserKeyLink(userId int, url *string) error
	GetRecipeKeyLink(recipeId int) (*string, error)
	SetRecipeKeyLink(recipeId int, url *string) error
}
