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
		usersRepo: usersRepo,
		recipesRepo: recipesRepo,
		categoriesRepo: categoriesRepo,
		apiKey: apiKey,
		firestore: firestore,
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

type SelectableStruct struct {
	Item     int  `json:"item"`
	Selected bool `json:"selected"`
}

func (s *FirebaseService) migrateFromFirebase(authData models.AuthData, firebaseUser models.FirebaseUser) error {
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
		var ingredients []SelectableStruct
		ingredients = firebaseRecipe["ingredients"].([]SelectableStruct)
		var steps []SelectableStruct
		steps = firebaseRecipe["cooking"].([]SelectableStruct)
		recipe := models.Recipe{
			Name:        firebaseRecipe["name"].(string),
			OwnerId:     userId,
			Favourite:   firebaseRecipe["favourite"].(bool),
			Servings:    firebaseRecipe["servings"].(int16),
			// Time     int16 `json:"time,omitempty"`
			Calories:    firebaseRecipe["calories"].(int16),
			Ingredients: ingredients,
			Cooking:     steps,
		}
		_, err := s.recipesRepo.CreateRecipe(recipe)
		if err != nil {
			logger.Error("failed to import recipe ", recipe.Name)
		}
	}
	return nil
}
