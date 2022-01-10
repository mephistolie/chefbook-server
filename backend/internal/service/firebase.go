package service

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repository"
	"github.com/mephistolie/chefbook-server/pkg/logger"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const FirebaseSignInEmailEndpoint = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword"

const FirestoreUsersCollection = "users"
const FirestoreRecipesCollection = "recipes"

type FirebaseService struct {
	usersRepo        repository.Users
	recipesRepo      repository.Recipes
	categoriesRepo   repository.Categories
	shoppingListRepo repository.ShoppingList
	apiKey           string
	firestore        firestore.Client
}

func NewFirebaseService(apiKey string, usersRepo repository.Users, recipesRepo repository.Recipes,
	categoriesRepo repository.Categories, shoppingListRepo repository.ShoppingList, firestore firestore.Client) *FirebaseService {
	return &FirebaseService{
		usersRepo:        usersRepo,
		recipesRepo:      recipesRepo,
		categoriesRepo:   categoriesRepo,
		shoppingListRepo: shoppingListRepo,
		apiKey:           apiKey,
		firestore:        firestore,
	}
}

func (s *FirebaseService) FirebaseSignIn(authData models.AuthData) (models.FirebaseUser, error) {
	route := fmt.Sprintf("%s?key=%s", FirebaseSignInEmailEndpoint, s.apiKey)

	jsonInput, err := json.Marshal(authData)
	if err != nil {
		return models.FirebaseUser{}, err
	}

	resp, err := http.Post(route, "application/json", bytes.NewBuffer(jsonInput))
	if err != nil {
		return models.FirebaseUser{}, err
	}
	var firebaseUser models.FirebaseUser
	err = json.NewDecoder(resp.Body).Decode(&firebaseUser)
	if err != nil {
		return models.FirebaseUser{}, err
	}
	return firebaseUser, nil
}

type MarkdownString struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

func (s *FirebaseService) migrateFromFirebase(authData models.AuthData, firebaseUser models.FirebaseUser) error {
	logger.Error(firebaseUser.LocalId)
	userSnapshot, err := s.firestore.Collection(FirestoreUsersCollection).Doc(firebaseUser.LocalId).Get(context.Background())
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

	username := userDoc["name"].(string)
	if len(username) > 0 {
		err = s.usersRepo.SetUserName(userId, username)
		if err != nil {
			logger.Warn("migration: error during setting name ", username)
		}
	}
	premium := userDoc["premium"].(bool)
	if premium {
		lifetimePremium := time.Date(3000, 1, 1, 0, 0, 0, 0, nil)
		err := s.usersRepo.SetPremiumDate(userId, lifetimePremium)
		if err != nil {
			logger.Warn("migration: error during activating premium")
		}
	}

	s.importFirebaseShoppingList(userId, userDoc)
	s.importFirebaseRecipes(userId, firebaseUser)
	return nil
}

func (s *FirebaseService) importFirebaseShoppingList(userId int, userDoc map[string]interface{}) {
	shoppingList := models.ShoppingList{
		Timestamp: time.Now(),
	}
	var firebaseShoppingList []string
	firebaseShoppingList = userDoc["shoppingList"].([]string)
	for _, firebasePurchase := range firebaseShoppingList {
		purchase := models.Purchase{
			Id:          uuid.NewString(),
			Item:        firebasePurchase,
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

func (s *FirebaseService) importFirebaseRecipes(userId int, firebaseUser models.FirebaseUser) {
	recipesSnapshot := s.firestore.Collection(FirestoreUsersCollection).Doc(firebaseUser.LocalId).Collection(FirestoreRecipesCollection).Documents(context.Background())
	firebaseRecipes, err := recipesSnapshot.GetAll()
	if err != nil {
		logger.Error("migration: error during get recipe list")
	} else {
		var categories []string
		var categoriesIds []int
		for _, firebaseRecipeSnapshot := range firebaseRecipes {
			firebaseRecipe := firebaseRecipeSnapshot.Data()

			var firebaseIngredients []interface{}
			firebaseIngredients = firebaseRecipe["ingredients"].([]interface{})
			var ingredients []MarkdownString
			for _, firebaseIngredient := range firebaseIngredients {
				mapIngredient := firebaseIngredient.(map[string]interface{})
				item := mapIngredient["item"].(string)
				selected := mapIngredient["selected"].(bool)
				stringType := "STRING"
				if selected {
					stringType = "HEADER"
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
				continue
			}

			var jsonCooking []byte
			if firebaseCooking, ok := firebaseRecipe["cooking"].([]interface{}); ok {
				var cooking []MarkdownString
				for _, firebaseStep := range firebaseCooking {
					mapStep := firebaseStep.(map[string]interface{})
					item := mapStep["item"].(string)
					selected := mapStep["selected"].(bool)
					stringType := "STRING"
					if selected {
						stringType = "HEADER"
					}
					step := MarkdownString{
						Text: item,
						Type: stringType,
					}
					cooking = append(cooking, step)
				}
				jsonCooking, err = json.Marshal(cooking)
				if err != nil {
					continue
				}
			} else {
				firebaseOldCooking, ok := firebaseRecipe["cooking"].([]string)
				if !ok {
					logger.Warn("migration: error during import steps of recipe")
					continue
				}
				var cooking []MarkdownString
				for _, firebaseStep := range firebaseOldCooking {
					step := MarkdownString{
						Text: firebaseStep,
						Type: "STRING",
					}
					cooking = append(cooking, step)
				}
				jsonCooking, err = json.Marshal(cooking)
				if err != nil {
					logger.Warn("migration: error during import steps of recipe")
					continue
				}
			}


			recipeTime := 0
			firebaseTime := firebaseRecipe["time"].(string)
			numberFilter := regexp.MustCompile("[0-9]+")
			timeSlice := numberFilter.FindAllString(firebaseTime, -1)
			timeSliceLength := len(timeSlice)
			if timeSliceLength > 0 {
				multiplier := 1
				if timeSliceLength == 1 && len(timeSlice[timeSliceLength-1]) == 1 {
					multiplier = 60
				}
				number, err := strconv.Atoi(timeSlice[timeSliceLength-1])
				if err == nil {
					recipeTime = recipeTime + number*multiplier
				}
			}
			if timeSliceLength > 1 {
				hours, err := strconv.Atoi(timeSlice[timeSliceLength-2])
				if err == nil {
					recipeTime = recipeTime + hours*60
				}
			}

			firebaseCategories := firebaseRecipe["categories"].([]string)
			var recipeCategoriesIds []int
			for _, category := range firebaseCategories {
				isAdded := false
				for i, addedCategory := range categories {
					if addedCategory == category {
						recipeCategoriesIds = append(recipeCategoriesIds, categoriesIds[i])
						isAdded = true
						break
					}
				}
				if !isAdded && len(category) > 0 {
					dbCategory := models.Category{
						Name: category,
						Cover: category[:1],
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

			recipe := models.Recipe{
				Name:        firebaseRecipe["name"].(string),
				OwnerId:     userId,
				Favourite:   firebaseRecipe["favourite"].(bool),
				Servings:    int16(firebaseRecipe["servings"].(int64)),
				Time:        int16(recipeTime),
				Calories:    int16(firebaseRecipe["calories"].(int64)),
				Ingredients: jsonIngredients,
				Cooking:     jsonCooking,
				Visibility:  "private",
			}
			recipeId, err := s.recipesRepo.CreateRecipe(recipe)
			if err != nil {
				continue
			}
			if err = s.recipesRepo.SetRecipeCategories(recipeCategoriesIds, recipeId, userId); err != nil {
				continue
			}
		}
	}
}