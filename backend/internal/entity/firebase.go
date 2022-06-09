package entity

import "time"

type FirebaseProfile struct {
	IdToken      string
	Email        string
	RefreshToken string
	ExpiresIn    string
	LocalId      string
	Registered   bool
}

type FirebaseProfileInfo struct {
	Username          *string
	CreationTimestamp *time.Time
	IsPremium         bool
}

type FirebaseRecipe struct {
	Recipe      RecipeInput
	Categories  []string
	IsFavourite bool
}

type FirebaseUserData struct {
	Profile      FirebaseProfileInfo
	Recipes      []FirebaseRecipe
	Categories   []CategoryInput
	ShoppingList ShoppingList
}
