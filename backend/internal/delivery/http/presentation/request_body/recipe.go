package request_body

import (
	"chefbook-server/internal/delivery/http/presentation/common_body"
	"chefbook-server/internal/entity"
	"chefbook-server/internal/entity/failure"
	"strings"
)

type RecipeInput struct {
	Name        string  `json:"name"`
	Visibility  string  `json:"visibility"`
	IsEncrypted bool    `json:"encrypted"`
	Language    string  `json:"language"`
	Description *string `json:"description"`
	Preview     *string `json:"preview"`

	Servings *int16 `json:"servings"`
	Time     *int16 `json:"time"`

	Calories       *int16                      `json:"calories"`
	Macronutrients *common_body.Macronutrients `json:"macronutrients"`

	Ingredients []common_body.IngredientItem `json:"ingredients"`
	Cooking     []common_body.CookingItem    `json:"cooking"`
}

func (r *RecipeInput) Validate() error {

	if len(r.Name) == 0 {
		return failure.EmptyRecipeName
	}

	if len(r.Name) > 100 {
		return failure.TooLongRecipeName
	}

	if len(r.Language) != 2 {
		r.Language = entity.CodeEnglish
	}
	r.Language = strings.ToLower(r.Language)

	r.Visibility = strings.ToLower(r.Visibility)
	if r.Visibility != entity.VisibilityShared && r.Visibility != entity.VisibilityPublic {
		r.Visibility = entity.VisibilityPrivate
	}

	if r.IsEncrypted && r.Visibility == entity.VisibilityPublic {
		return failure.InvalidEncryptionType
	}

	if r.Description != nil && len(*r.Description) > 1500 {
		return failure.TooLongRecipeDescription
	}

	if r.Time != nil && *r.Time <= 0 {
		r.Time = nil
	}

	if r.Calories != nil && *r.Calories <= 0 {
		r.Calories = nil
	}

	if r.Macronutrients != nil {
		if r.Macronutrients.Fats != nil && *r.Macronutrients.Fats <= 0 {
			r.Macronutrients.Fats = nil
		}
		if r.Macronutrients.Protein != nil && *r.Macronutrients.Protein <= 0 {
			r.Macronutrients.Protein = nil
		}
		if r.Macronutrients.Carbohydrates != nil && *r.Macronutrients.Carbohydrates <= 0 {
			r.Macronutrients.Carbohydrates = nil
		}
	}

	if len(r.Ingredients) == 0 {
		return failure.EmptyIngredients
	}

	if len(r.Cooking) == 0 {
		return failure.EmptyCooking
	}

	ingredientsEncrypted := false
	for _, ingredient := range r.Ingredients {
		if err := ingredient.Validate(); err != nil {
			return err
		}
		if ingredient.IsEncrypted() {
			ingredientsEncrypted = true
		}
	}
	if ingredientsEncrypted && len(r.Ingredients) != 1 {
		return failure.InvalidEncryptionType
	}

	cookingEncrypted := false
	for _, cookingItem := range r.Cooking {
		if err := cookingItem.Validate(); err != nil {
			return err
		}
		if cookingItem.IsEncrypted() {
			cookingEncrypted = true
		}
	}
	if cookingEncrypted && len(r.Cooking) != 1 {
		return failure.InvalidEncryptionType
	}

	if r.IsEncrypted && (!ingredientsEncrypted || !cookingEncrypted) || !r.IsEncrypted && (ingredientsEncrypted || cookingEncrypted) {
		return failure.InvalidEncryptionType
	}

	return nil
}

func (r *RecipeInput) Entity() entity.RecipeInput {
	macronutrients := common_body.Macronutrients{}
	if r.Macronutrients != nil {
		macronutrients = *r.Macronutrients
	}

	ingredients := make([]entity.IngredientItem, len(r.Ingredients))
	for i, ingredient := range r.Ingredients {
		ingredients[i] = ingredient.Entity()
	}

	cooking := make([]entity.CookingItem, len(r.Cooking))
	for i, cookingItem := range r.Cooking {
		cooking[i] = cookingItem.Entity()
	}

	return entity.RecipeInput{
		Name:        r.Name,
		Visibility:  r.Visibility,
		IsEncrypted: r.IsEncrypted,
		Language:    r.Language,
		Description: r.Description,
		Preview:     r.Preview,

		Servings: r.Servings,
		Time:     r.Time,

		Calories:       r.Calories,
		Macronutrients: macronutrients.Entity(),

		Ingredients: ingredients,
		Cooking:     cooking,
	}
}
