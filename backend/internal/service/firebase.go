package service

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/model"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"github.com/mephistolie/chefbook-server/pkg/logger"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

const FirebaseSignInEmailEndpoint = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword"

const firestoreUsersCollection = "users"
const firestoreRecipesCollection = "recipes"

const oldUserBroccoins = 500

type FirebaseService struct {
	usersRepo             repository.Auth
	profileRepo           repository.Profile
	recipesRepo           repository.RecipeCrud
	recipeInteractionRepo repository.RecipeInteraction
	categoriesRepo        repository.Categories
	shoppingListRepo      repository.ShoppingList
	apiKey                string
	firestore             firestore.Client
	auth                  auth.Client
}

func NewFirebaseService(apiKey string, usersRepo repository.Auth, profileRepo repository.Profile, recipesRepo repository.RecipeCrud,
	recipeInteractionRepo repository.RecipeInteraction, categoriesRepo repository.Categories, shoppingListRepo repository.ShoppingList,
	app firebase.App) *FirebaseService {
	firestoreClient, _ := app.Firestore(context.Background())
	firebaseAuth, _ := app.Auth(context.Background())
	return &FirebaseService{
		usersRepo:             usersRepo,
		profileRepo:           profileRepo,
		recipesRepo:           recipesRepo,
		recipeInteractionRepo: recipeInteractionRepo,
		categoriesRepo:        categoriesRepo,
		shoppingListRepo:      shoppingListRepo,
		apiKey:                apiKey,
		firestore:             *firestoreClient,
		auth:                  *firebaseAuth,
	}
}

func (s *FirebaseService) SignIn(authData model.AuthData) (model.FirebaseUser, error) {
	route := fmt.Sprintf("%s?key=%s", FirebaseSignInEmailEndpoint, s.apiKey)

	jsonInput, err := json.Marshal(authData)
	if err != nil {
		return model.FirebaseUser{}, err
	}

	resp, err := http.Post(route, "application/json", bytes.NewBuffer(jsonInput))
	if err != nil {
		return model.FirebaseUser{}, err
	}
	var firebaseUser model.FirebaseUser
	err = json.NewDecoder(resp.Body).Decode(&firebaseUser)
	if err != nil {
		return model.FirebaseUser{}, err
	}
	return firebaseUser, nil
}

type MarkdownString struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

func (s *FirebaseService) migrateFromFirebase(authData model.AuthData, firebaseUser model.FirebaseUser) error {
	userSnapshot, err := s.firestore.Collection(firestoreUsersCollection).Doc(firebaseUser.LocalId).Get(context.Background())
	userDoc := userSnapshot.Data()

	activationLink := uuid.New()
	userId, err := s.usersRepo.CreateUser(authData, activationLink)
	if err != nil {
		return err
	}
	err = s.usersRepo.ActivateUser(activationLink)
	if err != nil {
		logger.Warn("migration: error during activating link ", activationLink)
	}
	if user, err := s.auth.GetUser(context.Background(), firebaseUser.LocalId); err == nil {
		if err := s.profileRepo.SetProfileCreationDate(userId, time.UnixMilli(user.UserMetadata.CreationTimestamp)); err != nil {
			logger.Warn("migration: error during set profile creation date")
		}
	}
	err = s.profileRepo.IncreaseBroccoins(userId, oldUserBroccoins)
	if err != nil {
		logger.Warn("migration: error during adding broccoins to old user")
	}

	username, ok := userDoc["name"].(string)
	if ok && len(username) > 0 {
		err = s.profileRepo.SetUsername(userId, username)
		if err != nil {
			logger.Warn("migration: error during setting name ", username)
		}
	}
	premium, ok := userDoc["isPremium"].(bool)
	if ok && premium {
		lifetimePremium := time.Date(3000, 1, 1, 00, 00, 00, 00, time.UTC)
		err := s.profileRepo.SetPremiumDate(userId, lifetimePremium)
		if err != nil {
			logger.Warn("migration: error during activating premium")
		}
	}

	s.importFirebaseShoppingList(userId, userDoc)
	s.importFirebaseRecipes(userId, firebaseUser)
	return nil
}

func (s *FirebaseService) importFirebaseShoppingList(userId int, userDoc map[string]interface{}) {
	shoppingList := model.ShoppingList{
		Timestamp: time.Now(),
	}

	firebaseShoppingList, ok := userDoc["shoppingList"].([]interface{})
	if ok {
		for _, firebasePurchase := range firebaseShoppingList {
			purchase := model.Purchase{
				Id:          uuid.NewString(),
				Item:        firebasePurchase.(string),
				Multiplier:  1,
				IsPurchased: false,
			}
			shoppingList.Purchases = append(shoppingList.Purchases, purchase)
		}
		err := s.shoppingListRepo.SetShoppingList(shoppingList, userId)
		if err != nil {
			logger.Warn("migration: error during setting shopping list ")
		}
	}
}

func (s *FirebaseService) importFirebaseRecipes(userId int, firebaseUser model.FirebaseUser) {
	recipesSnapshot := s.firestore.Collection(firestoreUsersCollection).Doc(firebaseUser.LocalId).Collection(firestoreRecipesCollection).Documents(context.Background())
	firebaseRecipes, err := recipesSnapshot.GetAll()
	if err != nil {
		logger.Error("migration: error during get recipe list")
	} else {
		var categories []string
		var categoriesIds []int
		for _, firebaseRecipeSnapshot := range firebaseRecipes {
			firebaseRecipe := firebaseRecipeSnapshot.Data()

			name, ok := firebaseRecipe["name"].(string)
			if !ok {
				logger.Warn("migration: error during import name of recipe")
				continue
			}
			favourite, ok := firebaseRecipe["favourite"].(bool)
			if !ok {
				favourite = false
			}
			servings, ok := firebaseRecipe["servings"].(int64)
			if !ok {
				servings = 5
			}
			calories, ok := firebaseRecipe["calories"].(int64)
			if !ok {
				calories = 0
			}

			recipeTime := 60
			firebaseTime, ok := firebaseRecipe["time"].(string)
			if ok {
				recipeTime = parseTime(firebaseTime)
			}

			jsonIngredients, err := parseIngredients(firebaseRecipe)
			if err != nil {
				continue
			}

			jsonCooking, err := parseCooking(firebaseRecipe)
			if err != nil {
				continue
			}

			var recipeCategoriesIds []int
			firebaseCategories, ok := firebaseRecipe["categories"].([]interface{})
			if ok {
				for _, interfaceCategory := range firebaseCategories {
					category := interfaceCategory.(string)
					isAdded := false
					for i, addedCategory := range categories {
						if addedCategory == category {
							recipeCategoriesIds = append(recipeCategoriesIds, categoriesIds[i])
							isAdded = true
							break
						}
					}
					if !isAdded && len(category) > 0 {
						dbCategory := model.Category{
							Name:   category,
							Cover:  string([]rune(category)[0:1]),
							UserId: userId,
						}
						categoryId, err := s.categoriesRepo.AddCategory(dbCategory)
						if err == nil {
							categories = append(categories, category)
							categoriesIds = append(categoriesIds, categoryId)
							recipeCategoriesIds = append(recipeCategoriesIds, categoryId)
						}
					}
				}
			}

			recipe := model.Recipe{
				Name:        name,
				OwnerId:     userId,
				Language:    "en",
				Servings:    int16(servings),
				Time:        int16(recipeTime),
				Calories:    int16(calories),
				Ingredients: jsonIngredients,
				Cooking:     jsonCooking,
				Visibility:  "private",
			}
			recipeId, err := s.recipesRepo.CreateRecipe(recipe)
			if err != nil {
				continue
			}
			if err = s.recipeInteractionRepo.SetRecipeCategories(recipeCategoriesIds, recipeId, userId); err != nil {
				continue
			}
			if err = s.recipeInteractionRepo.SetRecipeFavourite(recipeId, userId, favourite); err != nil {
				continue
			}
		}
	}
}

func parseTime(timeString string) int {
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
	return minutes
}

func parseIngredients(firebaseRecipe map[string]interface{}) ([]byte, error) {
	var firebaseIngredients []interface{}
	firebaseIngredients, ok := firebaseRecipe["ingredients"].([]interface{})
	if !ok {
		logger.Warn("migration: error during import ingredients of recipe")
		return []byte{}, os.ErrInvalid
	}
	var ingredients []MarkdownString
	for _, firebaseIngredient := range firebaseIngredients {
		mapIngredient := firebaseIngredient.(map[string]interface{})
		var item string
		var selected bool
		nullableItem := mapIngredient["item"]
		nullableSelected := mapIngredient["selected"]
		if nullableItem == nil {
			nullableItem = mapIngredient["name"]
			nullableSelected = mapIngredient["section"]
		}
		item = nullableItem.(string)
		selected = nullableSelected.(bool)
		stringType := "INGREDIENT"
		if selected {
			stringType = "SECTION"
		}
		ingredient := MarkdownString{
			Text: item,
			Type: stringType,
		}
		ingredients = append(ingredients, ingredient)
	}
	jsonIngredients, err := json.Marshal(ingredients)
	if err != nil {
		logger.Warn("migration: error during import ingredients of recipe")
		return []byte{}, os.ErrInvalid
	}
	return jsonIngredients, nil
}

func parseCooking(firebaseRecipe map[string]interface{}) ([]byte, error) {
	var err error
	var jsonCooking []byte
	if firebaseCooking, ok := firebaseRecipe["cooking"].([]interface{}); ok {
		var cooking []MarkdownString
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

			stringType := "STEP"
			if selected {
				stringType = "SECTION"
			}
			step := MarkdownString{
				Text: item,
				Type: stringType,
			}
			cooking = append(cooking, step)
		}
		jsonCooking, err = json.Marshal(cooking)
		if err != nil {
			return []byte{}, os.ErrInvalid
		}
	}

	return jsonCooking, nil
}
