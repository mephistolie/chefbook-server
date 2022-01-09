package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mephistolie/chefbook-server/internal/models"
	"log"
	"net/http"
)

const FirebaseSignInEmailEndpoint = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword"

type FirebaseService struct {
	apiKey string
}

func NewFirebaseService(apiKey string) *FirebaseService {
	return &FirebaseService{
		apiKey: apiKey,
	}
}

func (s *FirebaseService) FirebaseSignIn(authData models.AuthData) (models.FirebaseUser, error) {
	route := fmt.Sprintf("%skey=%s", FirebaseSignInEmailEndpoint, s.apiKey)

	firebaseAuthData := map[interface{}]interface{} {
		"email": authData.Email,
		"password": authData.Password,
		"returnSecureToken": true,
	}
	jsonInput, err := json.Marshal(firebaseAuthData)
	if err != nil {
		return models.FirebaseUser{}, err
	}

	resp, err := http.Post(route, "application/json", bytes.NewBuffer(jsonInput))
	if err != nil {
		log.Fatal(err)
	}
	var firebaseUser models.FirebaseUser
	err = json.NewDecoder(resp.Body).Decode(&firebaseUser)
	if err != nil {
		return models.FirebaseUser{}, err
	}
	return firebaseUser, nil
}