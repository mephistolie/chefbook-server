package services

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/mephistolie/chefbook-server/internal/models"
	"github.com/mephistolie/chefbook-server/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	signingKey = ""
	tokenTTL = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repositories.Authorization
}

func NewAuthService(repo repositories.Authorization) *AuthService  {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	hashedPassword, err := generatePasswordHash(user.Password)
	if err != nil {
		return 0, nil
	}
	user.Password = hashedPassword
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(email, password string) (string, error)  {
	hashedPassword, err := generatePasswordHash(password)
	if err != nil {
		return "", nil
	}
	user, err := s.repo.GetUser(email, hashedPassword)
	if err != nil {
		return "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt: time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func generatePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hashedPassword), nil
}