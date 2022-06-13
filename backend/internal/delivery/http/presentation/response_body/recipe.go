package response_body

import (
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/common_body"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"time"
)

type Recipe struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	OwnerId     int     `json:"owner_id"`
	OwnerName   string  `json:"owner_name"`
	IsOwned     bool    `json:"owned"`
	IsSaved     bool    `json:"saved"`
	Likes       int16   `json:"likes"`
	Visibility  string  `json:"visibility"`
	IsEncrypted bool    `json:"encrypted"`
	Language    string  `json:"language"`
	Description *string `json:"description,omitempty"`
	Preview     *string `json:"preview,omitempty"`

	CreationTimestamp time.Time `json:"creation_timestamp"`
	UpdateTimestamp   time.Time `json:"update_timestamp"`

	Categories  *[]Category `json:"categories,omitempty"`
	IsFavourite bool        `json:"favourite"`
	IsLiked     bool        `json:"liked"`

	Servings *int16 `json:"servings,omitempty"`
	Time     *int16 `json:"time,omitempty"`

	Calories       *int16                      `json:"calories,omitempty"`
	Macronutrients *common_body.Macronutrients `json:"macronutrients,omitempty"`

	Ingredients []common_body.IngredientItem `json:"ingredients"`
	Cooking     []common_body.CookingItem    `json:"cooking"`
}

func NewRecipe(recipe entity.UserRecipe) Recipe {
	ingredients := make([]common_body.IngredientItem, len(recipe.Ingredients))
	for i, ingredient := range recipe.Ingredients {
		ingredients[i] = common_body.NewIngredientItem(ingredient)
	}

	cooking := make([]common_body.CookingItem, len(recipe.Cooking))
	for i, cookingItem := range recipe.Cooking {
		cooking[i] = common_body.NewCookingItem(cookingItem)
	}

	return Recipe{
		Id:          recipe.Id,
		Name:        recipe.Name,
		OwnerId:     recipe.OwnerId,
		OwnerName:   recipe.OwnerName,
		IsOwned:     recipe.IsOwned,
		IsSaved:     recipe.IsSaved,
		Likes:       recipe.Likes,
		Visibility:  recipe.Visibility,
		IsEncrypted: recipe.IsEncrypted,
		Language:    recipe.Language,
		Description: recipe.Description,
		Preview:     recipe.Preview,

		CreationTimestamp: recipe.CreationTimestamp.UTC(),
		UpdateTimestamp:   recipe.UpdateTimestamp.UTC(),

		Categories:  getRecipeCategories(&recipe.Categories),
		IsFavourite: recipe.IsFavourite,
		IsLiked:     recipe.IsLiked,

		Servings: recipe.Servings,
		Time:     recipe.Time,

		Calories:       recipe.Calories,
		Macronutrients: common_body.NewMacronutrients(recipe.Macronutrients),

		Ingredients: ingredients,
		Cooking:     cooking,
	}
}

type RecipeInfo struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	OwnerId     int     `json:"owner_id"`
	OwnerName   string  `json:"owner_name"`
	IsOwned     bool    `json:"owned"`
	IsSaved     bool    `json:"saved"`
	Likes       int16   `json:"likes"`
	Visibility  string  `json:"visibility"`
	IsEncrypted bool    `json:"encrypted"`
	Language    string  `json:"language"`
	Preview     *string `json:"preview,omitempty"`

	CreationTimestamp time.Time `json:"creation_timestamp"`
	UpdateTimestamp   time.Time `json:"update_timestamp"`

	Categories  *[]Category `json:"categories,omitempty"`
	IsFavourite bool        `json:"favourite"`
	IsLiked     bool        `json:"liked"`

	Servings *int16 `json:"servings,omitempty"`
	Time     *int16 `json:"time,omitempty"`

	Calories *int16 `json:"calories,omitempty"`
}

func NewRecipeInfo(recipe entity.RecipeInfo) RecipeInfo {
	return RecipeInfo{
		Id:          recipe.Id,
		Name:        recipe.Name,
		OwnerId:     recipe.OwnerId,
		OwnerName:   recipe.OwnerName,
		IsOwned:     recipe.IsOwned,
		IsSaved:     recipe.IsSaved,
		Likes:       recipe.Likes,
		Visibility:  recipe.Visibility,
		IsEncrypted: recipe.IsEncrypted,
		Language:    recipe.Language,
		Preview:     recipe.Preview,

		CreationTimestamp: recipe.CreationTimestamp.UTC(),
		UpdateTimestamp:   recipe.UpdateTimestamp.UTC(),

		Categories:  getRecipeCategories(&recipe.Categories),
		IsFavourite: recipe.IsFavourite,
		IsLiked:     recipe.IsLiked,

		Servings: recipe.Servings,
		Time:     recipe.Time,

		Calories: recipe.Calories,
	}
}

func getRecipeCategories(categories *[]entity.Category) *[]Category {
	var categoriesPointer *[]Category = nil
	if categories != nil && len(*categories) > 0 {
		responseCategories := make([]Category, len(*categories))
		for i, category := range *categories {
			responseCategories[i] = NewCategory(category)
		}
		categoriesPointer = &responseCategories
	}
	return categoriesPointer
}

func NewRecipes(entities []entity.RecipeInfo) []RecipeInfo {
	recipes := make([]RecipeInfo, len(entities))
	for i, recipe := range entities {
		recipes[i] = NewRecipeInfo(recipe)
	}
	return recipes
}
