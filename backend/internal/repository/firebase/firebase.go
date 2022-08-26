package firebase

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/internal/repository/firebase/dto"
	"github.com/mephistolie/chefbook-server/pkg/logger"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	firebaseSignInEmailEndpoint = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword"

	firestoreUsersCollection   = "users"
	firestoreRecipesCollection = "recipes"
)

type MigrationRepo struct {
	apiKey    string
	firestore firestore.Client
	auth      auth.Client
}

func NewMigrationRepo(app firebase.App, apiKey string) *MigrationRepo {
	firestoreClient, err := app.Firestore(context.Background())
	if err != nil {
		logger.Error(err)
	}
	firebaseAuth, err := app.Auth(context.Background())
	if err != nil {
		logger.Error(err)
	}
	return &MigrationRepo{
		apiKey:    apiKey,
		firestore: *firestoreClient,
		auth:      *firebaseAuth,
	}
}

func (r *MigrationRepo) GetProfile(credentials entity.Credentials) (entity.FirebaseProfile, error) {
	route := fmt.Sprintf("%s?key=%s", firebaseSignInEmailEndpoint, r.apiKey)

	firebaseCredentials := dto.NewFirebaseCredentials(credentials)
	jsonInput, err := json.Marshal(firebaseCredentials)
	if err != nil {
		return entity.FirebaseProfile{}, failure.InvalidCredentials
	}

	resp, err := http.Post(route, "application/json", bytes.NewBuffer(jsonInput))
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.Error(resp)
		return entity.FirebaseProfile{}, failure.InvalidCredentials
	}

	var firebaseUser dto.FirebaseProfile
	err = json.NewDecoder(resp.Body).Decode(&firebaseUser)
	if err != nil {
		return entity.FirebaseProfile{}, failure.UnableImportFirebaseProfile
	}

	return firebaseUser.Entity(), nil
}

func (r *MigrationRepo) GetProfileData(profile entity.FirebaseProfile) (entity.FirebaseUserData, error) {
	userSnapshot, err := r.firestore.Collection(firestoreUsersCollection).Doc(profile.LocalId).Get(context.Background())
	if err != nil {
		logger.Error("migration: ", err)
		return entity.FirebaseUserData{}, failure.UserNotFound
	}

	userDoc := userSnapshot.Data()

	user := r.getProfileInfo(profile, userDoc)
	recipes, categories := r.getRecipesAndCategories(profile)
	shoppingList := r.getShoppingList(userDoc)

	return entity.FirebaseUserData{
		Profile:      user,
		Recipes:      recipes,
		Categories:   categories,
		ShoppingList: shoppingList,
	}, nil
}

func (r *MigrationRepo) getProfileInfo(profile entity.FirebaseProfile, userDoc map[string]interface{}) entity.FirebaseProfileInfo {
	var firebaseProfileInfo entity.FirebaseProfileInfo

	if user, err := r.auth.GetUser(context.Background(), profile.LocalId); err == nil {
		timestamp := time.UnixMilli(user.UserMetadata.CreationTimestamp)
		firebaseProfileInfo.CreationTimestamp = &timestamp
	}

	username, ok := userDoc["name"].(string)
	if ok && len(username) > 0 {
		firebaseProfileInfo.Username = &username
	}

	premium, ok := userDoc["isPremium"].(bool)
	if ok && premium {
		firebaseProfileInfo.IsPremium = true
	}

	return firebaseProfileInfo
}

func (r *MigrationRepo) getRecipesAndCategories(firebaseUser entity.FirebaseProfile) ([]entity.FirebaseRecipe, []entity.CategoryInput) {
	var recipes []entity.FirebaseRecipe
	var categories []entity.CategoryInput

	recipesSnapshot := r.firestore.Collection(firestoreUsersCollection).Doc(firebaseUser.LocalId).Collection(firestoreRecipesCollection).Documents(context.Background())
	recipesDocs, err := recipesSnapshot.GetAll()
	if err != nil {
		return recipes, categories
	}

	for _, firebaseRecipeSnapshot := range recipesDocs {
		var ok bool
		recipe := entity.FirebaseRecipe{
			Recipe: entity.RecipeInput{
				Visibility:  entity.VisibilityPrivate,
				IsEncrypted: false,
				Language:    entity.CodeEnglish,
			},
		}
		recipeDoc := firebaseRecipeSnapshot.Data()

		recipe.Recipe.Name, ok = recipeDoc["name"].(string)
		if !ok {
			logger.Warn("migration: error during import name of recipe")
			continue
		}

		recipe.IsFavourite, _ = recipeDoc["favourite"].(bool)

		if servings, ok := recipeDoc["servings"].(int16); ok {
			recipe.Recipe.Servings = &servings
		}

		if firebaseTime, ok := recipeDoc["time"].(string); ok {
			recipeTime := parseTime(firebaseTime)
			recipe.Recipe.Time = &recipeTime
		}

		if calories, ok := recipeDoc["calories"].(int16); ok {
			recipe.Recipe.Calories = &calories
		}

		recipe.Recipe.Ingredients, err = parseIngredients(recipeDoc)
		if err != nil {
			logger.Error("parse ingredients")
			continue
		}

		recipe.Recipe.Cooking, err = parseCooking(recipeDoc)
		if err != nil {
			logger.Error("parse cooking")
			continue
		}

		firebaseCategories, ok := recipeDoc["categories"].([]interface{})
		if ok {
			for _, interfaceCategory := range firebaseCategories {
				category := interfaceCategory.(string)
				recipe.Categories = append(recipe.Categories, category)

				isAddedToAllCategories := false
				for _, addedCategory := range categories {
					if addedCategory.Name == category {
						isAddedToAllCategories = true
						break
					}
				}

				if !isAddedToAllCategories && len(category) > 0 {
					cover := string([]rune(category)[0:1])
					categoryEntity := entity.CategoryInput{
						Name:  category,
						Cover: &cover,
					}
					categories = append(categories, categoryEntity)
				}
			}
		}
		recipes = append(recipes, recipe)
	}

	return recipes, categories
}

func parseTime(timeString string) int16 {
	minutes := 0
	numberFilter := regexp.MustCompile("[0-9]+")
	timeSlice := numberFilter.FindAllString(timeString, -1)
	timeSliceLength := len(timeSlice)
	if timeSliceLength > 0 {
		multiplier := 1
		if timeSliceLength == 1 && len(timeSlice[timeSliceLength-1]) == 1 {
			multiplier = 60
		}
		number, err := strconv.Atoi(timeSlice[timeSliceLength-1])
		if err == nil {
			minutes = minutes + number*multiplier
		}
	}
	if timeSliceLength > 1 {
		hours, err := strconv.Atoi(timeSlice[timeSliceLength-2])
		if err == nil {
			minutes = minutes + hours*60
		}
	}
	return int16(minutes)
}

func parseIngredients(firebaseRecipe map[string]interface{}) ([]entity.IngredientItem, error) {
	firebaseIngredients, ok := firebaseRecipe["ingredients"].([]interface{})
	if !ok {
		logger.Warn("migration: error during import ingredients of recipe")
		return []entity.IngredientItem{}, failure.UnableImportFirebaseProfile
	}

	var ingredients []entity.IngredientItem
	for _, firebaseIngredient := range firebaseIngredients {
		var item string
		var selected bool
		mapIngredient := firebaseIngredient.(map[string]interface{})

		nullableItem := mapIngredient["item"]
		if nullableItem == nil {
			nullableItem = mapIngredient["name"]
			if nullableItem == nil {
				continue
			}
		}
		item = nullableItem.(string)

		nullableSelected := mapIngredient["selected"]
		if nullableSelected == nil {
			nullableSelected = mapIngredient["section"]
		}
		if nullableSelected != nil {
			selected = nullableSelected.(bool)
		}
		stringType := entity.TypeIngredient
		if selected {
			stringType = entity.TypeSection
		}

		ingredient := entity.IngredientItem{
			Id:   uuid.NewString(),
			Text: item,
			Type: stringType,
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

func parseCooking(firebaseRecipe map[string]interface{}) ([]entity.CookingItem, error) {
	firebaseCooking, ok := firebaseRecipe["cooking"].([]interface{})

	if !ok {
		logger.Warn("migration: error during import cooking of recipe")
		return []entity.CookingItem{}, failure.UnableImportFirebaseProfile
	}

	var cooking []entity.CookingItem
	for _, firebaseStep := range firebaseCooking {
		var item string
		var selected bool
		mapStep, ok := firebaseStep.(map[string]interface{})
		if ok {
			item = mapStep["item"].(string)
			selected = mapStep["selected"].(bool)
		} else {
			item = firebaseStep.(string)
		}

		stringType := entity.TypeStep
		if selected {
			stringType = entity.TypeSection
		}
		step := entity.CookingItem{
			Id:   uuid.NewString(),
			Text: item,
			Type: stringType,
		}
		cooking = append(cooking, step)
	}

	return cooking, nil
}

func (r *MigrationRepo) getShoppingList(userDoc map[string]interface{}) entity.ShoppingList {
	shoppingList := entity.ShoppingList{
		Timestamp: time.Now(),
	}

	firebaseShoppingList, ok := userDoc["shoppingList"].([]interface{})
	if ok {
		for _, firebasePurchase := range firebaseShoppingList {
			purchase := entity.Purchase{
				Id:          uuid.NewString(),
				Name:        firebasePurchase.(string),
				Multiplier:  1,
				IsPurchased: false,
			}
			shoppingList.Purchases = append(shoppingList.Purchases, purchase)
		}
	}

	return shoppingList
}
