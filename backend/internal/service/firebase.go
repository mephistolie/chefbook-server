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
			continue
		}

		var firebaseCooking []interface{}
		firebaseCooking = firebaseRecipe["cooking"].([]interface{})
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
		jsonCooking, err := json.Marshal(cooking)
		if err != nil {
			continue
		}

		time := 0
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
				time = time + number * multiplier
			}
		}
		if timeSliceLength > 1 {
			hours, err := strconv.Atoi(timeSlice[timeSliceLength-2])
			if err == nil {
				time = time + hours * 60
			}
		}

		recipe := models.Recipe{
			Name:        firebaseRecipe["name"].(string),
			OwnerId:     userId,
			Favourite:   firebaseRecipe["favourite"].(bool),
			Servings:    int16(firebaseRecipe["servings"].(int64)),
			Time:        int16(time),
			Calories:    int16(firebaseRecipe["calories"].(int64)),
			Ingredients: jsonIngredients,
			Cooking:     jsonCooking,
			Visibility:  "private",
		}
		_, err = s.recipesRepo.CreateRecipe(recipe)
		if err != nil {
			continue
		}
	}
	return nil
}
