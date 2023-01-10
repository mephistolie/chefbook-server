package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type TokenManager interface {
	NewJWT(userId uuid.UUID, ttl time.Duration) (string, time.Time, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(userId uuid.UUID, ttl time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().Add(ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: expiresAt.Unix(),
		Subject:   userId.String(),
	})

	signedString, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", time.Now(), err
	}
	return signedString, expiresAt, nil
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
