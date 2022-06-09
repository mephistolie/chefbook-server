package entity

import (
	"time"
)

const (
	VisibilityPrivate = "private"
	VisibilityShared  = "shared"
	VisibilityPublic  = "public"

	CodeEnglish = "en"
)

type Recipe struct {
	Id          int
	Name        string
	OwnerId     int
	Likes       int16
	Visibility  string
	IsEncrypted bool
	Language    string
	Description *string
	Preview     *string

	CreationTimestamp time.Time
	UpdateTimestamp   time.Time

	Servings *int16
	Time     *int16

	Calories       *int16
	Macronutrients Macronutrients

	Ingredients []IngredientItem
	Cooking     []CookingItem
}

type UserRecipe struct {
	Id          int
	Name        string
	OwnerId     int
	OwnerName   string
	Owned       bool
	Likes       int16
	Visibility  string
	IsEncrypted bool
	Language    string
	Description *string
	Preview     *string

	CreationTimestamp time.Time
	UpdateTimestamp   time.Time

	Categories  []Category
	IsFavourite bool
	IsLiked     bool

	Servings *int16
	Time     *int16

	Calories       *int16
	Macronutrients Macronutrients

	Ingredients []IngredientItem
	Cooking     []CookingItem
}

type RecipeInfo struct {
	Id          int
	Name        string
	OwnerId     int
	OwnerName   string
	Owned       bool
	Likes       int16
	Visibility  string
	IsEncrypted bool
	Language    string
	Preview     *string

	CreationTimestamp time.Time
	UpdateTimestamp   time.Time

	Categories  []Category
	IsFavourite bool
	IsLiked     bool

	Servings *int16
	Time     *int16

	Calories *int16
}

type RecipeInput struct {
	Name        string
	Visibility  string
	IsEncrypted bool
	Language    string
	Description *string
	Preview     *string

	Servings *int16
	Time     *int16

	Calories       *int16
	Macronutrients Macronutrients

	Ingredients []IngredientItem
	Cooking     []CookingItem
}
