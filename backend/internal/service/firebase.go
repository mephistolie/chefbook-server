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
)

const FirebaseSignInEmailEndpoint = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword"

const FirestoreUsersCollection = "users"
const FirestoreRecipesCollection = "recipes"

type FirebaseService struct {
	usersRepo      repository.Users
	recipesRepo    repository.Recipes
	categoriesRepo repository.Categories
	apiKey         string
	firestore      firestore.Client
}

func NewFirebaseService(apiKey string, usersRepo repository.Users, recipesRepo repository.Recipes,
	categoriesRepo repository.Categories, firestore firestore.Client) *FirebaseService {
	return &FirebaseService{
		usersRepo:      usersRepo,
		recipesRepo:    recipesRepo,
		categoriesRepo: categoriesRepo,
		apiKey:         apiKey,
		firestore:      firestore,
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
		return err
	}
	username := userDoc["name"].(string)
	err = s.usersRepo.SetUserName(userId, username)
	if err != nil {
		return err
	}

	recipesSnapshot := s.firestore.Collection(FirestoreUsersCollection).Doc(firebaseUser.LocalId).Collection(FirestoreRecipesCollection).Documents(context.Background())
	firebaseRecipes, err := recipesSnapshot.GetAll()
	if err != nil {
		return err
	}
	for _, firebaseRecipeSnapshot := range firebaseRecipes {
		firebaseRecipe := firebaseRecipeSnapshot.Data()

		var firebaseIngredients []interface{}
		firebaseIngredients = firebaseRecipe["ingredients"].([]interface{})
		var ingredients []MarkdownString
		for _, firebaseIngredient := range firebaseIngredients {
			mapIngredient := firebaseIngredient.(map[string]interface{})
			item := mapIngredient["item"].(string)
			selected := mapIngredient["selected"].(bool)
			stringType := "string"
			if selected {
				stringType = "header"
			}
			ingredient := MarkdownString {
				Text: item,
				Type: stringType,
			}
			ingredients = append(ingredients, ingredient)
		}
		jsonIngredients, err := json.Marshal(ingredients)
		if err != nil {
			continue
		}

		var firebaseCooking []interface{}
		firebaseCooking = firebaseRecipe["cooking"].([]interface{})
		var cooking []MarkdownString
		for _, firebaseStep := range firebaseCooking {
			mapStep := firebaseStep.(map[string]interface{})
			item := mapStep["item"].(string)
			selected := mapStep["selected"].(bool)
			stringType := "string"
			if selected {
				stringType = "header"
			}
			step := MarkdownString {
				Text: item,
				Type: stringType,
			}
			cooking = append(cooking, step)
		}
		jsonCooking, err := json.Marshal(cooking)
		if err != nil {
			continue
		}

		recipe := models.Recipe{
			Name:      firebaseRecipe["name"].(string),
			OwnerId:   userId,
			Favourite: firebaseRecipe["favourite"].(bool),
			Servings:  firebaseRecipe["servings"].(int16),
			// Time     int16 `json:"time,omitempty"`
			Calories:    firebaseRecipe["calories"].(int16),
			Ingredients: jsonIngredients,
			Cooking:     jsonCooking,
		}
		_, err = s.recipesRepo.CreateRecipe(recipe)
		if err != nil {
			continue
		}
	}
	return nil
}
